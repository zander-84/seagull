package worker

import (
	"context"
	"errors"
	"github.com/zander-84/seagull/contract"
	"github.com/zander-84/seagull/think"
	"github.com/zander-84/seagull/tool"
	"log"
	"runtime"
	"sync"
	"time"
)

type processor struct {
	waiter sync.WaitGroup
}

func NewProcessor() contract.Processor {
	out := new(processor)
	return out
}

func (p *processor) Go(handler func(), w *sync.WaitGroup) {
	p.waiter.Add(1)
	if w != nil {
		w.Add(1)
	}
	go func() {
		defer p.waiter.Done()
		if w != nil {
			defer w.Done()
		}
		defer func() {
			if rErr := recover(); rErr != nil {
				buf := make([]byte, 64<<10)
				n := runtime.Stack(buf, false)
				buf = buf[:n]
				log.Printf("Printf err: %v \n", rErr)
				log.Println(string(buf))

			}
		}()

		handler()
	}()
}

func (p *processor) GoTimeout(handler func() error, timeout time.Duration) (err error) {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	fin := make(chan error, 1)
	go func() {
		defer func() {
			if rErr := recover(); rErr != nil {
				buf := make([]byte, 64<<10)
				n := runtime.Stack(buf, false)
				buf = buf[:n]
				log.Printf("Printf err: %v \n", rErr)
				log.Println(string(buf))

				err = errors.New(string(buf))
			}
		}()
		fin <- handler()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case e := <-fin:
		return e
	}
}

func (p *processor) GoListenCtx(contexts []context.Context, handler func() (any, error)) (any, error) {
	var out = make(chan interface{}, 1)
	var outErr = make(chan error, 2)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 64<<10)
				n := runtime.Stack(buf, false)
				buf = buf[:n]
				log.Printf("Printf err: %v \n", r)
				log.Println(string(buf))

				outErr <- errors.New(string(buf))
			}
		}()
		tmp, err := handler()
		if err != nil {
			outErr <- err
			return
		}
		out <- tmp
	}()

	var newContexts = make([]context.Context, 0)
	for _, v := range contexts {
		if v.Done() != nil {
			newContexts = append(newContexts, v)
		}
	}
	newContexts2Len := len(newContexts)
	if newContexts2Len == 0 {
		for {
			select {
			case o := <-out:
				return o, nil
			case e := <-outErr:
				return nil, e
			}
		}
	} else if newContexts2Len == 1 {
		select {
		case o := <-out:
			return o, nil
		case e := <-outErr:
			return nil, e
		case <-newContexts[0].Done():
			return nil, errors.New("main process exit")
		}
	} else if newContexts2Len == 2 {
		select {
		case o := <-out:
			return o, nil
		case e := <-outErr:
			return nil, e
		case <-newContexts[0].Done():
			return nil, errors.New("main process exit")
		case <-newContexts[1].Done():
			return nil, errors.New("main process exit")
		}
	} else {
		quitCh, cancel2 := think.DoneCtxChan(newContexts...)
		defer cancel2()
		for {
			select {
			case o := <-out:
				return o, nil
			case e := <-outErr:
				return nil, e
			case <-quitCh:
				return nil, errors.New("main process exit")
			}
		}
	}

}

func (p *processor) Wait(duration time.Duration) error {
	return tool.ExitWithTimeout(duration, func() error {
		p.waiter.Wait()
		return nil
	})
}
