package limit

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

const (
	luaIncrScript = `
local key = KEYS[1]
local incrVal = tonumber(ARGV[1])
local ttl = tonumber(ARGV[2])
local minVal = tonumber(ARGV[3])
local maxVal = tonumber(ARGV[4])
local allow  = 0

local currentVal = redis.call("get", key)
local isFirst = false
if  currentVal == false then
	currentVal = 0
	isFirst = true
end

if not isFirst then
	currentVal = tonumber(currentVal)
end

if currentVal + incrVal < minVal then
	allow = 0
elseif currentVal + incrVal > maxVal then
	allow = 0
else 
	allow = 1
	currentVal = redis.call("incrby", key, ARGV[1])
end

if isFirst then
	if ttl > 0 then
		redis.call("pexpire", key, ARGV[2])
	end
end

local lessTtl = redis.call("pttl", key)
-- 自动更新剩余时间 [剩余时间>自定义时间 | 剩余时间永久 自定义时间大于0 | 剩余时间大于0 自定义时间等于永久]
if not isFirst and (lessTtl > ttl or (lessTtl < 0  and ttl > 0) or (lessTtl > 0 and ttl < 1)) then
	 if ttl < 1 then
			redis.call("PERSIST", key)
	else 
			redis.call("pexpire", key, ARGV[2])
	end
	lessTtl = redis.call("pttl", key)
end
return {currentVal, lessTtl, allow }
`

	luaPeekScript = `
local key = KEYS[1]
local v = redis.call("get", key)
if v == false then
	return {0,0,1}
end
local ttl = redis.call("pttl", key)
return {tonumber(v), ttl,1}
`
)

type limitContainer struct {
	engine     *redis.Client
	luaPeekSHA string
	luaIncrSHA string
	luaLoaded  uint32
	luaMutex   sync.RWMutex
}

func NewLimitContainer(engine *redis.Client) (contract.Limiter, error) {
	l := &limitContainer{engine: engine}
	err := l.preloadLuaScripts(context.Background())
	return l, err
}

func (l *limitContainer) loadLuaScripts(ctx context.Context) error {
	l.luaMutex.Lock()
	defer l.luaMutex.Unlock()

	if atomic.LoadUint32(&l.luaLoaded) != 0 {
		return nil
	}

	luaPeekSHA, err := l.engine.ScriptLoad(ctx, luaPeekScript).Result()
	if err != nil {
		return err
	}
	luaIncrSHA, err := l.engine.ScriptLoad(ctx, luaIncrScript).Result()
	if err != nil {
		return err
	}

	l.luaPeekSHA = luaPeekSHA
	l.luaIncrSHA = luaIncrSHA

	atomic.StoreUint32(&l.luaLoaded, 1)

	return nil
}

func (l *limitContainer) preloadLuaScripts(ctx context.Context) error {
	if atomic.LoadUint32(&l.luaLoaded) == 0 {
		return l.loadLuaScripts(ctx)
	}
	return nil
}

func (l *limitContainer) reloadLuaScripts(ctx context.Context) error {
	atomic.StoreUint32(&l.luaLoaded, 0)
	return l.loadLuaScripts(ctx)
}

func (l *limitContainer) Get(ctx context.Context, key string) (int64, error) {
	count, _, _, err := l.get(ctx, key)
	return count, err
}

func (l *limitContainer) get(ctx context.Context, key string) (int64, time.Duration, bool, error) {
	cmd := l.evalSHA(ctx, l.getLuaPeekSHA, []string{key})
	count, ttl, ok, err := parseCountAndTTL(cmd)
	return count, ttl, ok, err
}

func (l *limitContainer) Allow(ctx context.Context, key string, incrVal int64, expires time.Duration, minVal, maxVal int64) (int64, time.Duration, bool, error) {
	cmd := l.evalSHA(ctx, l.getLuaIncrSHA, []string{key}, incrVal, expires.Milliseconds(), minVal, maxVal)
	count, ttl, ok, err := parseCountAndTTL(cmd)
	return count, ttl, ok, err
}

func (l *limitContainer) evalSHA(ctx context.Context, getSha func() string,
	keys []string, args ...interface{}) *redis.Cmd {

	cmd := l.engine.EvalSha(ctx, getSha(), keys, args...)
	err := cmd.Err()
	if err == nil || !isLuaScriptGone(err) {
		return cmd
	}

	err = l.reloadLuaScripts(ctx)
	if err != nil {
		cmd = redis.NewCmd(ctx)
		cmd.SetErr(err)
		return cmd
	}

	return l.engine.EvalSha(ctx, getSha(), keys, args...)
}

func (l *limitContainer) getLuaPeekSHA() string {
	l.luaMutex.RLock()
	defer l.luaMutex.RUnlock()
	return l.luaPeekSHA
}
func (l *limitContainer) getLuaIncrSHA() string {
	l.luaMutex.RLock()
	defer l.luaMutex.RUnlock()
	return l.luaIncrSHA
}

func isLuaScriptGone(err error) bool {
	return strings.HasPrefix(err.Error(), "NOSCRIPT")
}

func parseCountAndTTL(cmd *redis.Cmd) (int64, time.Duration, bool, error) {
	result, err := cmd.Result()
	if err != nil {
		return 0, 0, false, err
	}

	fields, ok := result.([]interface{})
	if !ok || len(fields) != 3 {
		return 0, 0, false, errors.New("three elements in result were expected")
	}

	currentVal, ok1 := fields[0].(int64)
	ttl, ok2 := fields[1].(int64)
	okTag, ok3 := fields[2].(int64)
	if !ok1 || !ok2 || !ok3 {
		return 0, 0, false, errors.New("type of the count and/or ttl should be number and/or act res should be number")
	}
	var actRes = false
	if okTag == 1 {
		actRes = true
	}

	return currentVal, time.Millisecond * time.Duration(ttl), actRes, nil
}
