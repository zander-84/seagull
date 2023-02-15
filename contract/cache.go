package contract

import (
	"context"
	"github.com/zander-84/seagull/contract/def"
	"time"
)

type Cache interface {
	Exists(ctx context.Context, keys ...def.K) (bool, error)
	Get(ctx context.Context, key def.K, recPtr any) error
	Set(ctx context.Context, key def.K, value any, expires time.Duration) error
	SetNX(ctx context.Context, key def.K, value any, expires time.Duration) (bool, error)

	Delete(ctx context.Context, keys ...def.K) error

	// DelayDelete  用于延迟双删
	DelayDelete(ctx context.Context, delay time.Duration, keys ...def.K) error
	GetOrSet(ctx context.Context, key def.K, recPtr any, expires time.Duration, f func(key def.K) (value any, err error)) error

	//BatchGetOrSet  返回的：value map[string]any string表示缓存的key
	BatchGetOrSet(ctx context.Context, keys []def.K, recPtr any, expires time.Duration, f func(missIds []def.K) (value map[string]any, err error)) error
	FlushDB(ctx context.Context) error
	Ping(ctx context.Context) error
}
