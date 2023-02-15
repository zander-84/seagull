package goredis

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"sync"
	"sync/atomic"
)

type HashCtl struct {
	mutex sync.RWMutex

	luaAddLoaded uint32
	luaAddScript string
	luaAddSHA    string

	luaSubLoaded uint32
	luaSubScript string
	luaSubSHA    string
	engine       *redis.Client
}

// -5:库存操作失败
// -4:代表库存传进来的值是负数（非法值）
// -3:库存未初始化
// -2:库存不足
// -1:库存为0大于等于0:剩余库存（扣减之后剩余的库存)
func newHashCtl(engine *redis.Client) (*HashCtl, error) {
	out := new(HashCtl)
	out.engine = engine
	out.luaSubScript = `
local key = KEYS[1]
local filed = KEYS[2]
local count = tonumber(ARGV[1])
local minCount = tonumber(ARGV[2])

local exist = redis.call("hexists",key,filed)
if exist < 1 then 
	return {exist,-3}
end

local nowRet = redis.call("hget",key,filed)
if nowRet - count < minCount then 
	return {tonumber(nowRet),-2}
end
local ret = redis.call("hincrby", key,filed, -ARGV[1])
return {ret, -1}
`
	out.luaAddScript = `
local key = KEYS[1]
local filed = KEYS[2]
local count = tonumber(ARGV[1])
local maxCount = tonumber(ARGV[2])

local exist = redis.call("hexists",key,filed)
if exist < 1 then 
	redis.call("hset",key,filed,0)
end

local nowRet = redis.call("hget",key,filed)
if nowRet + count > maxCount then 
	return {tonumber(nowRet),-2}
end
local ret = redis.call("hincrby", key,filed, ARGV[1])
return {ret, -1}
`
	if err := out.preloadSubLuaScript(context.Background()); err != nil {
		return out, err
	}
	if err := out.preloadAddLuaScript(context.Background()); err != nil {
		return out, err
	}
	return out, nil
}

func (h *HashCtl) preloadSubLuaScript(ctx context.Context) error {
	if atomic.LoadUint32(&h.luaSubLoaded) == 0 {
		return h.loadSubLuaScript(ctx)
	}
	return nil
}
func (h *HashCtl) preloadAddLuaScript(ctx context.Context) error {
	if atomic.LoadUint32(&h.luaSubLoaded) == 0 {
		return h.loadAddLuaScript(ctx)
	}
	return nil
}
func (h *HashCtl) loadAddLuaScript(ctx context.Context) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if atomic.LoadUint32(&h.luaAddLoaded) != 0 {
		return nil
	}
	sha, err := h.engine.ScriptLoad(ctx, h.luaAddScript).Result()
	if err != nil {
		return err
	}
	h.luaAddSHA = sha
	atomic.StoreUint32(&h.luaAddLoaded, 1)

	return nil
}
func (h *HashCtl) loadSubLuaScript(ctx context.Context) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if atomic.LoadUint32(&h.luaSubLoaded) != 0 {
		return nil
	}
	sha, err := h.engine.ScriptLoad(ctx, h.luaSubScript).Result()
	if err != nil {
		return err
	}
	h.luaSubSHA = sha
	atomic.StoreUint32(&h.luaSubLoaded, 1)

	return nil
}

func (h *HashCtl) reloadSubLuaScripts(ctx context.Context) error {
	atomic.StoreUint32(&h.luaSubLoaded, 0)
	return h.loadSubLuaScript(ctx)
}
func (h *HashCtl) reloadAddLuaScripts(ctx context.Context) error {
	atomic.StoreUint32(&h.luaAddLoaded, 0)
	return h.loadAddLuaScript(ctx)
}

func (h *HashCtl) getLuaSubSHA() string {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return h.luaSubSHA
}
func (h *HashCtl) getLuaAddSHA() string {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return h.luaAddSHA
}

func (h *HashCtl) parseIncr(cmd *redis.Cmd) (int64, int64, error) {
	result, err := cmd.Result()
	if err != nil {
		return 0, 0, err
	}
	fields, ok := result.([]interface{})
	if !ok || len(fields) != 2 {
		return 0, 0, errors.New("two elements in result were expected")
	}

	cnt, ok1 := fields[0].(int64)
	res, ok2 := fields[1].(int64)
	if !ok1 || !ok2 {
		return 0, 0, errors.New("type of the count and/or res should be number")
	}

	return cnt, res, nil
}

// SubInventory 减库存，库存不能低于最小库存
// 0:当前库存 1: tag, 2：{true:成功，false：库存不足} 3:系统错误
func (h *HashCtl) SubInventory(ctx context.Context, key string, field string, val uint64, minInventory uint64) (int64, int64, bool, error) {
	cmd := evalSHA(ctx, h.engine, h.getLuaSubSHA, h.reloadSubLuaScripts, []string{key, field}, val, minInventory)
	cnt, tag, err := h.parseIncr(cmd)
	if err != nil {
		return 0, tag, false, err
	}
	if tag == -2 || tag == -3 {
		return cnt, tag, false, nil
	}
	if tag == -1 {
		return cnt, tag, true, nil
	}
	return cnt, tag, false, nil
}

func (h *HashCtl) IncrUnsafeInventory(ctx context.Context, key string, field string, val int64) (int64, error) {
	cmd := h.engine.HIncrBy(ctx, key, field, val)
	if cmd.Err() != nil {
		return 0, cmd.Err()
	}
	return cmd.Val(), nil
}

// AddSoldAndCompareMaxInventory ，添加售出,并且小于最大库存，否则添加失败
// 0:当前已出售数量 1: tag, 2：{true:成功，false：库存不足} 3:系统错误
func (h *HashCtl) AddSoldAndCompareMaxInventory(ctx context.Context, key string, field string, inventory uint64, maxInventory uint64) (int64, int64, bool, error) {
	cmd := evalSHA(ctx, h.engine, h.getLuaAddSHA, h.reloadAddLuaScripts, []string{key, field}, inventory, maxInventory)
	cnt, tag, err := h.parseIncr(cmd)
	if err != nil {
		return 0, tag, false, err
	}
	if tag == -2 || tag == -3 {
		return cnt, tag, false, nil
	}
	if tag == -1 {
		return cnt, tag, true, nil
	}
	return cnt, tag, false, nil
}

func (h *HashCtl) SetInventory(ctx context.Context, key string, field string, val int64) (int64, error) {
	cmd := h.engine.HSet(ctx, key, field, val)
	if cmd.Err() != nil {
		return 0, cmd.Err()
	}
	return val, nil
}

func (h *HashCtl) GetInventory(ctx context.Context, key string, field string) (int64, error) {
	cmd := h.engine.HGet(ctx, key, field)
	if cmd.Err() != nil {
		return 0, cmd.Err()
	}
	return cmd.Int64()
}

func (h *HashCtl) DelField(ctx context.Context, key string, field string) error {
	cmd := h.engine.HDel(ctx, key, field)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}
func (h *HashCtl) Set(ctx context.Context, key string, in interface{}) error {
	var tempMap map[string]interface{}

	if in2, ok := in.(map[string]interface{}); ok {
		tempMap = in2
	} else {
		temp, err := json.Marshal(in)
		if err != nil {
			return err
		}
		tempMap = make(map[string]interface{}, 0)
		if err := json.Unmarshal(temp, &tempMap); err != nil {
			return err
		}
	}

	dst := make([]interface{}, 0)
	for k, v := range tempMap {
		dst = append(dst, k, v)
	}
	cmd := h.engine.HSet(ctx, key, dst...)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

// GetByScan 需要设置tag
//
//	type testData struct {
//			Name string `redis:"Name"`
//			Age  int64  `redis:"Age"`
//		}
func (h *HashCtl) GetByScan(ctx context.Context, key string, dst interface{}) error {
	cmd := h.engine.HGetAll(ctx, key)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	if err := cmd.Scan(dst); err != nil {
		return err
	}
	return nil
}

func (h *HashCtl) GetByMap(ctx context.Context, key string) (map[string]string, error) {
	cmd := h.engine.HGetAll(ctx, key)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	return cmd.Val(), nil
}

func (h *HashCtl) FieldExist(ctx context.Context, key string, field string) (bool, error) {
	cmd := h.engine.HExists(ctx, key, field)
	if cmd.Err() != nil {
		return false, cmd.Err()
	}
	return cmd.Result()
}

func (h *HashCtl) GetFieldString(ctx context.Context, key string, field string) (string, error) {
	cmd := h.engine.HGet(ctx, key, field)

	if cmd.Err() != nil {
		return "", cmd.Err()
	}
	return cmd.Result()
}

func (h *HashCtl) GetFieldInt64(ctx context.Context, key string, field string) (int64, error) {
	cmd := h.engine.HGet(ctx, key, field)
	if cmd.Err() != nil {
		return 0, cmd.Err()
	}
	return cmd.Int64()
}

// GetFieldCtl 自己去解析
func (h *HashCtl) GetFieldCtl(ctx context.Context, key string, field string) (*redis.StringCmd, error) {
	cmd := h.engine.HGet(ctx, key, field)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	return cmd, nil
}
