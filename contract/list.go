package contract

import "context"

// List 适用于table缓存，可用于库存，以及一些临时的动态数据
type List interface {
	Len(ctx context.Context, key string) (int64, error)

	RightPop(ctx context.Context, key string, data any, checkData func(data any) (bool, error)) error
	RightPush(ctx context.Context, key string, data any) error
	RightBatchPush(ctx context.Context, key string, locker Locker, num int, purchaser func(i int) (any, error)) error

	LeftPop(ctx context.Context, key string, data any, checkData func(data any) (bool, error)) error
	LeftPush(ctx context.Context, key string, data any) error
	LeftBatchPush(ctx context.Context, key string, locker Locker, num int, purchaser func(i int) (any, error)) error

	Range(ctx context.Context, key string, ptrSliceData interface{}, offset, limit int64) error

	Del(ctx context.Context, key string) error
}
