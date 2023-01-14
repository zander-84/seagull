package lock

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/zander-84/seagull/contract"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const leaserRedisScript = `
local key = KEYS[1]
local id = ARGV[1]
local ttl = tonumber(ARGV[2])
local currentVal = redis.call("get", key)
local exitTag  = 0

if currentVal==id then
	exitTag = 1
	if ttl > 0 then
		redis.call("pexpire", key, ARGV[2])
	end
end
return {exitTag}
`

type redisLeaser struct {
	leaserInterval time.Duration
	leaser         map[string]chan struct{}
	leaserLocker   sync.RWMutex
	luaLeaserSHA   string
	luaLoaded      uint32
	luaMutex       sync.RWMutex
	processor      contract.Processor
	engine         *redis.Client
}

func newRedisLeaser(engine *redis.Client, processor contract.Processor, leaserInterval time.Duration) (*redisLeaser, error) {
	out := new(redisLeaser)
	out.engine = engine
	out.leaserInterval = leaserInterval
	out.leaser = make(map[string]chan struct{}, 0)
	out.processor = processor
	err := out.loadLuaScripts(context.Background())
	return out, err
}

func (r *redisLeaser) key(key string, id string) string {
	return key + ":" + id
}

func (r *redisLeaser) addLeaser(key string, id string) <-chan struct{} {
	r.leaserLocker.Lock()
	defer r.leaserLocker.Unlock()
	realKey := r.key(key, id)
	if _, ok := r.leaser[realKey]; !ok {
		r.leaser[realKey] = make(chan struct{}, 0)
	}

	return r.leaser[realKey]
}

func (r *redisLeaser) getLeasers() []string {
	r.leaserLocker.Lock()
	defer r.leaserLocker.Unlock()
	out := make([]string, 0, len(r.leaser))
	for k, _ := range r.leaser {
		out = append(out, k)
	}
	return out
}

func (r *redisLeaser) exitAndDelLeaser(key string, id string) {
	r.leaserLocker.Lock()
	defer r.leaserLocker.Unlock()

	realKey := r.key(key, id)
	leaser, ok := r.leaser[realKey]
	if ok {
		close(leaser)
		delete(r.leaser, realKey)
	}
}

func (r *redisLeaser) lease(ctx context.Context, key string, id string, expiration time.Duration) {
	leaserChan := r.addLeaser(key, id)
	r.processor.Go(func() {
		for {
			select {
			case <-leaserChan:
				return
			case <-time.After(r.leaserInterval):
				// 续租
				ok, err := r.doLua(ctx, key, id, expiration)
				if err != nil {
					ok, err = r.doLua(ctx, key, id, expiration)
					if err != nil {
						time.Sleep(time.Second / 2)
						ok, err = r.doLua(ctx, key, id, expiration)
					}
				}
				if !ok || err != nil {
					return
				}
			}
		}
	}, nil)
}

func (r *redisLeaser) release(ctx context.Context, key string, identify string) error {
	r.exitAndDelLeaser(key, identify)
	return nil
}

func (r *redisLeaser) doLua(ctx context.Context, key string, id string, expiration time.Duration) (bool, error) {
	cmd := r.evalSHA(ctx, r.getLuaLeaserSHA, []string{key}, id, expiration.Milliseconds())
	ok, err := parseLeaserLuaCmd(cmd)
	return ok, err
}
func (r *redisLeaser) evalSHA(ctx context.Context, getSha func() string,
	keys []string, args ...interface{}) *redis.Cmd {

	cmd := r.engine.EvalSha(ctx, getSha(), keys, args...)
	err := cmd.Err()
	if err == nil || !isLuaScriptGone(err) {
		return cmd
	}

	err = r.reloadLuaScripts(ctx)
	if err != nil {
		cmd = redis.NewCmd(ctx)
		cmd.SetErr(err)
		return cmd
	}

	return r.engine.EvalSha(ctx, getSha(), keys, args...)
}

func (r *redisLeaser) reloadLuaScripts(ctx context.Context) error {
	atomic.StoreUint32(&r.luaLoaded, 0)
	return r.loadLuaScripts(ctx)
}

func (r *redisLeaser) getLuaLeaserSHA() string {
	r.luaMutex.RLock()
	defer r.luaMutex.RUnlock()
	return r.luaLeaserSHA
}
func (r *redisLeaser) loadLuaScripts(ctx context.Context) error {
	r.luaMutex.Lock()
	defer r.luaMutex.Unlock()

	if atomic.LoadUint32(&r.luaLoaded) != 0 {
		return nil
	}

	luaLeaserSHA, err := r.engine.ScriptLoad(ctx, leaserRedisScript).Result()
	if err != nil {
		return err
	}
	r.luaLeaserSHA = luaLeaserSHA
	atomic.StoreUint32(&r.luaLoaded, 1)

	return nil
}

func isLuaScriptGone(err error) bool {
	return strings.HasPrefix(err.Error(), "NOSCRIPT")
}

func parseLeaserLuaCmd(cmd *redis.Cmd) (bool, error) {
	result, err := cmd.Result()
	if err != nil {
		return false, err
	}

	fields, ok := result.([]interface{})
	if !ok || len(fields) != 1 {
		return false, errors.New("one element in result were expected")
	}

	okTag, ok3 := fields[0].(int64)
	if !ok3 {
		return false, errors.New("type of the count and/or ttl should be number and/or act res should be number")
	}
	var actRes = false
	if okTag == 1 {
		actRes = true
	}

	return actRes, nil
}
