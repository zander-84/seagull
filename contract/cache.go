package contract

import (
	"context"
	"time"
)

type CacheKey interface {
	Key() string
	Raw() any
}
type cacheKey struct {
	key string
	raw any
}

func (c *cacheKey) Key() string {
	return c.key
}
func (c *cacheKey) Raw() any {
	return c.raw
}
func NewCacheKey(key string, raws ...any) CacheKey {
	out := new(cacheKey)
	out.key = key
	if len(raws) > 0 {
		out.raw = raws[0]
	}
	return out
}

type Cache interface {
	Exists(ctx context.Context, keys ...CacheKey) (bool, error)
	Get(ctx context.Context, key CacheKey, recPtr any) error
	Set(ctx context.Context, key CacheKey, value any, expires time.Duration) error
	SetNX(ctx context.Context, key CacheKey, value any, expires time.Duration) (bool, error)

	Delete(ctx context.Context, keys ...CacheKey) error

	// DelayDelete  用于延迟双删
	DelayDelete(ctx context.Context, delay time.Duration, keys ...CacheKey) error
	GetOrSet(ctx context.Context, key CacheKey, recPtr any, expires time.Duration, f func(key CacheKey) (value any, err error)) error

	BatchGetOrSet(ctx context.Context, keys []CacheKey, recPtr any, expires time.Duration, f func(missIds []CacheKey) (value map[CacheKey]any, err error)) error
	FlushDB(ctx context.Context) error
	Ping(ctx context.Context) error
}
