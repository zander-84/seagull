package middleware

import (
	"context"
	"github.com/zander-84/seagull/endpoint"
	"github.com/zander-84/seagull/think"
	"log"
	"runtime"
)

func Recover() endpoint.Middleware {
	return func(next endpoint.HandlerFunc) endpoint.HandlerFunc {
		return func(ctx context.Context, request interface{}) (out interface{}, err error) {
			defer func() {
				if rErr := recover(); rErr != nil {
					buf := make([]byte, 64<<10)
					n := runtime.Stack(buf, false)
					buf = buf[:n]
					log.Printf("Printf err: %v \n", rErr)
					log.Println(string(buf))
					err = think.ErrSystemSpace("from recover")
				}
			}()
			return next(ctx, request)
		}
	}
}
