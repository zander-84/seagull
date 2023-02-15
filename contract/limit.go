package contract

import (
	"context"
	"time"
)

type Limiter interface {
	Get(ctx context.Context, key string) (val int64, err error)

	// Allow
	// incrVal: 自增的值，可以是负数
	// expires：过期时间（如果系统的过期时间大于给于的值，就会修改系统的过期时间
	// minVal: 最小值  currentVal+incrVal 必须大于最小值
	// maxVal： 最大值  currentVal+incrVal 必须小于最大值
	//
	// val: 当前值
	// ttl: 剩余时间
	// ok: 操作是否成功
	// err: 错误
	Allow(ctx context.Context, key string, incrVal int64, expires time.Duration, minVal, maxVal int64) (val int64, ttl time.Duration, ok bool, err error)
}
