package tool

import (
	"context"
	"errors"
	"log"
	"runtime"
	"time"
)

func ExitWithTimeout(duration time.Duration, job func() error) error {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	return ExitWithContext(ctx, job)
}
func ExitWithContext(ctx context.Context, job func() error) error {
	done := ctx.Done()
	if done == nil {
		return job()
	}

	fin := make(chan error, 1)
	var err error

	go func() {
		func() {
			defer func() {
				if recoverErr := recover(); recoverErr != nil {
					buf := make([]byte, 64<<10)
					n := runtime.Stack(buf, false)
					buf = buf[:n]
					log.Printf("Printf err: %v \n", recoverErr)
					log.Println(string(buf))
					err = errors.New("exit panic")
				}
			}()
			err = job()
		}()
		fin <- err
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case finErr := <-fin:
		return finErr
	}
}
