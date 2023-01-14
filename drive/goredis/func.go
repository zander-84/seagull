package goredis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
	"strings"
	"time"
)

func (r *Rdb) TryLockWithTimeout(ctx context.Context, key string, identify string, duration time.Duration) (bool, error) {
	return r.engine.SetNX(ctx, key, identify, duration).Result()
}

func (r *Rdb) TryLockWithWaiting(ctx context.Context, key string, identify string, duration time.Duration, waitTime int) (bool, error) {
	for i := 0; i < waitTime; i++ {
		ok, err := r.engine.SetNX(ctx, key, identify, duration).Result()
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
		time.Sleep(time.Second)
	}
	return false, nil
}

func (r *Rdb) ReleaseLock(ctx context.Context, key string, identify string) error {
	data, err := r.GetString(ctx, key)
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}

	if data == identify {
		if err := r.Dels(ctx, key); err != nil {
			if err == redis.Nil {
				return nil
			}
			return err
		}
	}
	return err
}

func (r *Rdb) SetString(ctx context.Context, key string, str string, expires time.Duration) (err error) {
	return r.engine.Set(ctx, key, str, expires).Err()
}

func (r *Rdb) GetString(ctx context.Context, key string) (string, error) {
	return r.engine.Get(ctx, key).Result()
}

func (r *Rdb) GetInt64(ctx context.Context, key string) (int64, error) {
	stringResult, err := r.engine.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	if out, err1 := strconv.ParseInt(stringResult, 10, 64); err1 == nil {
		return out, nil
	}
	return 0, nil
}

func (r *Rdb) GetBytes(ctx context.Context, key string) ([]byte, error) {
	return r.engine.Get(ctx, key).Bytes()
}

func (r *Rdb) Dels(ctx context.Context, keys ...string) (err error) {
	return r.engine.Del(ctx, keys...).Err()
}

func (r *Rdb) HashHelper() *HashCtl {
	return r.hashCtl
}

func evalSHA(ctx context.Context, engine *redis.Client, getSha func() string, reloadLuaScript func(ctx2 context.Context) error,
	keys []string, args ...interface{}) *redis.Cmd {

	cmd := engine.EvalSha(ctx, getSha(), keys, args...)
	err := cmd.Err()
	if err == nil || !strings.HasPrefix(err.Error(), "NOSCRIPT") {
		return cmd
	}

	err = reloadLuaScript(ctx)
	if err != nil {
		cmd = redis.NewCmd(ctx)
		cmd.SetErr(err)
		return cmd
	}

	return engine.EvalSha(ctx, getSha(), keys, args...)
}
