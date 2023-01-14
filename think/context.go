package think

import (
	"context"
	"sync"
)

func HasDoneCtx(ctx ...context.Context) bool {
	for _, ctx1 := range ctx {
		if ok := hasDoneCtx(ctx1); ok {
			return true
		}
	}
	return false
}

func hasDoneCtx(ctx context.Context) bool {
	if ctx.Done() == nil {
		return false
	}
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

func DoneCtxChan(contexts ...context.Context) (quitCh <-chan struct{}, cancel func()) {
	var newContexts = make([]context.Context, 0)
	for _, v := range contexts {
		if v.Done() != nil {
			newContexts = append(newContexts, v)
		}
	}

	doneCh := make(chan struct{}, len(newContexts))
	quit := make(chan struct{})
	quitOnce := sync.Once{}
	cancel = func() {
		quitOnce.Do(func() {
			close(quit)
		})
	}
	for _, ctx := range newContexts {
		if ctx.Done() != nil {
			go func(ctx context.Context) {
				select {
				case <-ctx.Done():
					doneCh <- struct{}{}
					cancel()
				case <-quit:
					return
				}
			}(ctx)
		}

	}
	return doneCh, cancel
}
