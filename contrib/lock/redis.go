package lock

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/zander-84/seagull/contract"
	"time"
)

type redisLock struct {
	engine      *redis.Client
	unique      contract.Unique
	redisLeaser *redisLeaser
}
type redisLocked struct {
	key    string
	id     string
	engine *redisLock
}

func (r *redisLocked) Release(ctx context.Context) error {
	return r.engine.Release(ctx, r.key, r.id)
}

func (r *redisLocked) GetID() string {
	return r.id
}

func newRedisLocked(engine *redisLock, key string, id string) contract.Locked {
	out := new(redisLocked)
	out.key = key
	out.id = id
	out.engine = engine
	return out
}
func NewRedisLocker(engine *redis.Client, unique contract.Unique, processor contract.Processor) (contract.Locker, error) {
	var err error
	out := &redisLock{
		engine: engine,
		unique: unique,
	}
	out.redisLeaser, err = newRedisLeaser(engine, processor)
	return out, err
}

func (r *redisLock) Lock(ctx context.Context, key string, minExpiration time.Duration) (contract.Locked, error) {
	identify := r.unique.ID()
	ok, err := r.engine.SetNX(ctx, key, identify, minExpiration).Result()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, contract.LockFailed
	}
	if err == nil && ok {
		// 续租  Release 时候释放
		if minExpiration > 0 {
			r.redisLeaser.lease(ctx, key, identify, minExpiration)
		}
	}
	return newRedisLocked(r, key, identify), nil
}

func (r *redisLock) Release(ctx context.Context, key string, identify string) error {
	_ = r.redisLeaser.release(ctx, key, identify)

	data, err := r.engine.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}

	if data == identify || identify == "" {
		if err := r.engine.Del(ctx, key).Err(); err != nil {
			if err == redis.Nil {
				return nil
			}
			return err
		}
	}
	return err
}
