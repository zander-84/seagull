package cache

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/golang/groupcache/singleflight"
	"github.com/zander-84/seagull/contract"
	"github.com/zander-84/seagull/contract/def"
	"github.com/zander-84/seagull/think"
	"reflect"
	"runtime"
	"sync"
	"time"
)

type redisCache struct {
	engine       *redis.Client
	singleFlight singleflight.Group
	lock         sync.RWMutex
	codec        contract.Codec
	processor    contract.Processor
	try          int
}

func NewRedisCache(engine *redis.Client, codec contract.Codec, processor contract.Processor, try int) contract.Cache {
	return &redisCache{engine: engine, codec: codec, processor: processor, try: try}
}
func (r *redisCache) Exists(ctx context.Context, keys ...def.K) (bool, error) {
	var ok bool
	var err error
	for i := 0; i <= r.try; i++ {
		ok, err = r.exists(ctx, keys...)
		if err == nil {
			return ok, err
		}
	}

	return ok, err
}
func (r *redisCache) exists(ctx context.Context, keys ...def.K) (bool, error) {
	var rdsKey = make([]string, 0, len(keys))
	for _, v := range keys {
		rdsKey = append(rdsKey, v.Key)
	}
	cmd := r.engine.Exists(ctx, rdsKey...)
	if cmd.Err() != nil {
		return false, cmd.Err()
	}

	if cmd.Val() == int64(len(keys)) {
		return true, nil
	} else {
		return false, nil
	}
}

func (r *redisCache) Get(ctx context.Context, key def.K, recPtr any) (err error) {
	for i := 0; i <= r.try; i++ {
		err = r.get(ctx, key, recPtr)
		if err == nil {
			return nil
		} else {
			if think.IsErrNotFound(err) {
				return err
			}
		}
	}

	return err
}

func (r *redisCache) get(ctx context.Context, key def.K, recPtr any) (err error) {
	defer func() {
		if recoverErr := recover(); recoverErr != nil {
			buf := make([]byte, 64<<10)
			n := runtime.Stack(buf, false)
			buf = buf[:n]
			err = errors.New(string(buf))
		}
	}()

	data, err := r.singleFlight.Do(key.Key, func() (interface{}, error) {
		b, err := r.engine.Get(ctx, key.Key).Bytes()
		if err == redis.Nil {
			return nil, think.RecordNotFound
		}
		if err != nil {
			b, err = r.engine.Get(ctx, key.Key).Bytes()
			if err == redis.Nil {
				return nil, think.RecordNotFound
			}
			if err != nil {
				return nil, err
			}
		}
		err = r.codec.Unmarshal(b, recPtr)
		if err != nil {
			return nil, err
		}
		return recPtr, nil
	})

	if err != nil {
		return err
	}

	if data != nil && recPtr != data {
		reflect.ValueOf(recPtr).Elem().Set(getValue(data))
	}

	return err
}

func (r *redisCache) Set(ctx context.Context, key def.K, value any, expires time.Duration) error {
	var err error
	for i := 0; i <= r.try; i++ {
		err = r.set(ctx, key, value, expires)
		if err == nil {
			return nil
		}
	}

	return err
}

func (r *redisCache) set(ctx context.Context, key def.K, value any, expires time.Duration) error {
	b, err := r.codec.Marshal(value)
	if err != nil {
		return err
	}

	return r.engine.Set(ctx, key.Key, b, expires).Err()
}

func (r *redisCache) SetNX(ctx context.Context, key def.K, value any, expires time.Duration) (bool, error) {
	var ok bool
	var err error
	for i := 0; i <= r.try; i++ {
		ok, err = r.setNX(ctx, key, value, expires)
		if err == nil {
			return ok, err
		}
	}

	return ok, err
}

func (r *redisCache) setNX(ctx context.Context, key def.K, value any, expires time.Duration) (bool, error) {
	b, err := r.codec.Marshal(value)
	if err != nil {
		return false, err
	}

	return r.engine.SetNX(ctx, key.Key, b, expires).Result()
}

func (r *redisCache) Delete(ctx context.Context, keys ...def.K) error {
	var err error
	for i := 0; i <= r.try; i++ {
		err = r.delete(ctx, keys...)
		if err == nil {
			return err
		}
	}

	return err
}

func (r *redisCache) delete(ctx context.Context, keys ...def.K) error {
	var rdsKey = make([]string, 0, len(keys))
	for _, v := range keys {
		rdsKey = append(rdsKey, v.Key)
	}
	return r.engine.Del(ctx, rdsKey...).Err()
}

func (r *redisCache) DelayDelete(ctx context.Context, delay time.Duration, keys ...def.K) error {
	r.processor.Go(func() {
		select {
		case <-time.After(delay):
			if err := r.Delete(context.Background(), keys...); err != nil {
				err = r.Delete(context.Background(), keys...)
			}
		}
	}, nil)

	return nil
}

func (r *redisCache) GetOrSet(ctx context.Context, key def.K, recPtr any, expires time.Duration, f func(key def.K) (value any, err error)) (err error) {
	defer func() {
		if recoverErr := recover(); recoverErr != nil {
			buf := make([]byte, 64<<10)
			n := runtime.Stack(buf, false)
			buf = buf[:n]
			err = errors.New(string(buf))
		}
	}()

	if err := r.Get(ctx, key, recPtr); err != nil {
		if !think.IsErrNotFound(err) {
			return err
		}

		if fv, fe := r.singleFlight.Do(key.Key, func() (interface{}, error) {
			return f(key)
		}); fe != nil {
			return fe
		} else {
			// 不允许覆盖
			if _, err := r.SetNX(ctx, key, fv, expires); err != nil {
				return err
			}
			reflect.ValueOf(recPtr).Elem().Set(getValue(fv))
		}
	}

	return nil
}

func (r *redisCache) BatchGetOrSet(ctx context.Context, keys []def.K, recPtr any, expires time.Duration, f func(missIds []def.K) (value map[string]any, err error)) error {
	if reflect.ValueOf(recPtr).Elem().Type().Kind() != reflect.Slice {
		return errors.New("data  must be slice ptr")
	}
	reflectValue := reflect.ValueOf(recPtr).Elem()
	var missIds = make([]def.K, 0)
	for _, id := range keys {
		tmp := reflect.New(reflectValue.Type().Elem())
		err := r.Get(ctx, id, tmp.Interface())
		if err != nil {
			if think.IsErrNotFound(err) {
				missIds = append(missIds, id)
			} else {
				return err
			}
		} else {
			reflectValue.Set(reflect.Append(reflectValue, _getValue(tmp.Elem())))
		}
	}
	if len(missIds) < 1 {
		return nil
	}
	missVal, err := f(missIds)
	if err != nil {
		return err
	}
	for key, val := range missVal {
		if _, err := r.SetNX(ctx, def.K{Key: key}, val, expires); err != nil {
			return err
		}
		reflectValue.Set(reflect.Append(reflectValue, getValue(val)))
	}
	return err
}

func (r *redisCache) FlushDB(ctx context.Context) error {
	return r.engine.FlushDB(ctx).Err()
}

func (r *redisCache) Ping(ctx context.Context) error {
	return r.engine.Ping(ctx).Err()
}
