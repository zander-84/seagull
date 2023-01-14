package cache

import (
	"context"
	"errors"
	"github.com/golang/groupcache/singleflight"
	"github.com/patrickmn/go-cache"
	"github.com/zander-84/seagull/contract"
	"github.com/zander-84/seagull/think"
	"reflect"
	"runtime"
	"sync"
	"time"
)

type memory struct {
	engine       *cache.Cache
	singleFlight singleflight.Group
	lock         sync.RWMutex
	processor    contract.Processor
}

func NewMemoryCache(engine *cache.Cache, processor contract.Processor) contract.Cache {
	return &memory{engine: engine, processor: processor}
}

func (m *memory) Exists(ctx context.Context, keys ...contract.CacheKey) (bool, error) {
	for _, key := range keys {
		if _, ok := m.engine.Get(key.Key()); !ok {
			return false, nil
		}
	}
	return true, nil
}

func (m *memory) Get(ctx context.Context, key contract.CacheKey, recPtr interface{}) error {
	value, ok := m.engine.Get(key.Key())
	if !ok {
		return think.RecordNotFound
	}
	return m.unmarshal(value, recPtr)
}

func (m *memory) Set(ctx context.Context, key contract.CacheKey, value interface{}, expires time.Duration) error {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if expires == 0 {
		expires = -1
	}
	m.engine.Set(key.Key(), value, expires)
	return nil
}

func (m *memory) SetNX(ctx context.Context, key contract.CacheKey, value interface{}, expires time.Duration) (bool, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	ok, _ := m.Exists(ctx, key)
	if ok {
		return false, nil
	}

	if expires == 0 {
		expires = -1
	}

	m.engine.Set(key.Key(), value, expires)
	return true, nil
}

func (m *memory) Delete(ctx context.Context, key ...contract.CacheKey) error {
	for _, _key := range key {
		m.engine.Delete(_key.Key())
	}
	return nil
}
func (m *memory) DelayDelete(ctx context.Context, delay time.Duration, keys ...contract.CacheKey) error {
	m.processor.Go(func() {
		select {
		case <-time.After(delay):
			if err := m.Delete(context.Background(), keys...); err != nil {
				err = m.Delete(context.Background(), keys...)
			}
		}
	}, nil)

	return nil
}
func (m *memory) GetOrSet(ctx context.Context, key contract.CacheKey, recPtr interface{}, expires time.Duration, f func(key contract.CacheKey) (value any, err error)) (err error) {
	defer func() {
		if recoverErr := recover(); recoverErr != nil {
			err = think.ErrType("类型错误")
		}
	}()
	err = m.Get(ctx, key, recPtr)
	if err == nil {
		return nil
	}
	if !think.IsErrNotFound(err) {
		return err
	}

	fv, fe := m.singleFlight.Do(key.Key(), func() (any, error) {
		return f(key)
	})
	if fe != nil {
		err = fe
		return err
	}
	if _, err := m.SetNX(ctx, key, fv, expires); err != nil {
		return err
	}
	err = m.unmarshal(fv, recPtr)
	return err

}

func (m *memory) BatchGetOrSet(ctx context.Context, ids []contract.CacheKey, recPtr interface{}, expires time.Duration, f func(missIds []contract.CacheKey) (res map[contract.CacheKey]any, err error)) error {
	if reflect.ValueOf(recPtr).Elem().Type().Kind() != reflect.Slice {
		return errors.New("data  must be slice ptr")
	}
	reflectValue := reflect.ValueOf(recPtr).Elem()
	var missIds = make([]contract.CacheKey, 0)
	for _, id := range ids {
		tmp := reflect.New(reflectValue.Type().Elem())
		err := m.Get(ctx, id, tmp.Interface())
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
		if _, err := m.SetNX(ctx, key, val, expires); err != nil {
			return err
		}
		reflectValue.Set(reflect.Append(reflectValue, getValue(val)))
	}
	return err
}

func (m *memory) FlushDB(ctx context.Context) error {
	m.engine.Flush()
	return nil
}

func (m *memory) Ping(ctx context.Context) error {
	return nil
}
func (m *memory) unmarshal(form interface{}, toPtr interface{}) (err error) {
	defer func() {
		if rerr := recover(); rerr != nil {
			buf := make([]byte, 64<<10)
			n := runtime.Stack(buf, false)
			buf = buf[:n]
			err = think.ErrSystemSpace(string(buf))
		}
	}()

	reflect.ValueOf(toPtr).Elem().Set(getValue(form))
	return err
}
