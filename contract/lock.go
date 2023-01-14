package contract

import (
	"context"
	"errors"
	"time"
)

var LockFailed = errors.New("locking failed")

type Locked interface {
	Release(ctx context.Context) error
	GetID() string
}
type Locker interface {
	// Lock 建议用ctx控制时间
	// minExpiration 最小到期时间，用于续租异常兜底
	Lock(ctx context.Context, key string, minExpiration time.Duration) (locked Locked, err error)

	Release(ctx context.Context, key string, identify string) error
}
