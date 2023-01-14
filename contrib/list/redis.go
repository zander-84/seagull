package list

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/zander-84/seagull/contract"
	"github.com/zander-84/seagull/think"
	"reflect"
	"time"
)

type redisList struct {
	codec  contract.Codec
	engine *redis.Client
}

func NewRedisList(engine *redis.Client, codec contract.Codec) contract.List {
	return &redisList{
		engine: engine,
		codec:  codec,
	}
}
func (r *redisList) Len(ctx context.Context, key string) (int64, error) {
	cmd := r.engine.LLen(ctx, key)
	if err := cmd.Err(); err != nil {
		return 0, err
	}
	return cmd.Val(), nil
}

func (r *redisList) RightPop(ctx context.Context, key string, data any, checkData func(data any) (bool, error)) error {
	cmd := r.engine.RPop(ctx, key)
	if err := cmd.Err(); err != nil {
		if err == redis.Nil {
			return think.ErrRecordNotFound(err.Error())
		}
		return err
	}
	dataBytes, _ := cmd.Bytes()
	if len(dataBytes) == 0 {
		return think.ErrRecordNotFound("队列返回数据空")
	}
	if err := r.codec.Unmarshal(dataBytes, data); err != nil {
		return err
	}

	if checkData != nil {
		if ok, err := checkData(data); err != nil {
			return err
		} else if !ok {
			return r.RightPop(ctx, key, data, checkData)
		}
	}
	return nil
}

func (r *redisList) RightPush(ctx context.Context, key string, data any) error {
	dataBytes, err := r.codec.Marshal(data)
	if err != nil {
		return err
	}
	cmd := r.engine.RPush(ctx, key, dataBytes)
	if err = cmd.Err(); err != nil {
		return err
	}
	return nil
}

func (r *redisList) RightBatchPush(ctx context.Context, key string, locker contract.Locker, num int, purchaser func(i int) (any, error)) error {
	locked, err := locker.Lock(ctx, key+":lock", time.Minute)
	if err != nil {
		if err == contract.LockFailed {
			return nil
		}
		return err
	}

	defer locked.Release(ctx)

	for i := 0; i < num; i++ {
		data, err := purchaser(i)
		if err != nil {
			return err
		}
		if err := r.RightPush(ctx, key, data); err != nil {
			return err
		}
	}
	return nil
}
func (r *redisList) LeftPop(ctx context.Context, key string, data any, checkData func(data any) (bool, error)) error {
	cmd := r.engine.LPop(ctx, key)
	if err := cmd.Err(); err != nil {
		return err
	}
	dataBytes, _ := cmd.Bytes()
	if len(dataBytes) == 0 {
		return think.ErrRecordNotFound("队列返回数据空")
	}
	if err := r.codec.Unmarshal(dataBytes, data); err != nil {
		return err
	}
	if checkData != nil {
		if ok, err := checkData(data); err != nil {
			return err
		} else if !ok {
			return r.LeftPop(ctx, key, data, checkData)
		}
	}
	return nil
}

func (r *redisList) LeftPush(ctx context.Context, key string, data any) error {
	dataBytes, err := r.codec.Marshal(data)
	if err != nil {
		return err
	}
	cmd := r.engine.LPush(ctx, key, dataBytes)
	if err = cmd.Err(); err != nil {
		return err
	}
	return nil
}

func (r *redisList) LeftBatchPush(ctx context.Context, key string, locker contract.Locker, num int, purchaser func(i int) (any, error)) error {
	locked, err := locker.Lock(ctx, key+":lock", time.Minute)
	if err != nil {
		if err == contract.LockFailed {
			return nil
		}
		return err
	}

	defer locked.Release(ctx)

	for i := 0; i < num; i++ {
		data, err := purchaser(i)
		if err != nil {
			return err
		}
		if err := r.LeftPush(ctx, key, data); err != nil {
			return err
		}
	}
	return nil
}
func (r *redisList) Range(ctx context.Context, key string, ptrSliceData interface{}, offset, limit int64) error {
	if reflect.ValueOf(ptrSliceData).Type().Kind() != reflect.Ptr {
		return errors.New("data  must be ptr type")
	}
	if reflect.ValueOf(ptrSliceData).Elem().Type().Kind() != reflect.Slice {
		return errors.New("data  must be slice ptr")
	}
	start := offset
	var stop int64
	if offset >= 0 {
		stop = limit + offset - 1
	} else {
		stop = offset + limit + 1 //-5 2  -5 -4
		if stop >= 0 {
			stop = -1
		}
	}
	cmd := r.engine.LRange(ctx, key, start, stop)
	if err := cmd.Err(); err != nil {
		return err
	}

	reflectValue := reflect.ValueOf(ptrSliceData).Elem()

	for _, v := range cmd.Val() {
		tmp := reflect.New(reflectValue.Type().Elem())
		if err := r.codec.Unmarshal([]byte(v), tmp.Interface()); err != nil {
			return err
		}
		reflectValue.Set(reflect.Append(reflectValue, tmp.Elem()))
	}
	return nil
}

func (r *redisList) Del(ctx context.Context, key string) error {
	return r.engine.Del(ctx, key).Err()
}
