package contract

import (
	"context"
	"sync"
	"time"
)

type Processor interface {
	//Go 异步执行
	Go(handler func(), w *sync.WaitGroup)

	// GoTimeout 同步+超时
	GoTimeout(handler func() error, timeout time.Duration) error

	// GoListenCtx 同步 同任意context退出而退出
	GoListenCtx(contexts []context.Context, handler func() (any, error)) (any, error)

	//Wait 等待
	Wait(duration time.Duration) error
}
