package middleware

import (
	"context"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"github.com/ulule/limiter/v3/drivers/store/redis"
	"github.com/zander-84/seagull/endpoint"
	"github.com/zander-84/seagull/think"
)

var limiterErr = think.ErrTooManyRequests("访问过于频繁")
var limiterSystemErr = think.ErrAlert("限流器配置错误")

func Limiter(l *limiter.Limiter, idFunc func(ctx context.Context) (ip string, err error)) endpoint.Middleware {
	return func(next endpoint.HandlerFunc) endpoint.HandlerFunc {
		return func(ctx context.Context, request interface{}) (out interface{}, err error) {

			if ip, err := idFunc(ctx); err == nil {
				limiterContext, err := l.Get(ctx, ip)
				if err != nil {
					return nil, limiterSystemErr
				}

				if limiterContext.Reached {
					return nil, limiterErr
				}
			}

			return next(ctx, request)
		}
	}
}

// NewMemoryLimiter 内存限流
// * 5 reqs/second: "5-S"
// * 10 reqs/minute: "10-M"
// * 1000 reqs/hour: "1000-H"
// * 2000 reqs/day: "2000-D"
func NewMemoryLimiter(formatted string) *limiter.Limiter {
	_rate, err := limiter.NewRateFromFormatted(formatted)
	if err != nil {
		panic(err)
	}
	store := memory.NewStore()
	return limiter.New(store, _rate)
}

func NewRedisLimiter(formatted string, rds redis.Client) *limiter.Limiter {
	_rate, err := limiter.NewRateFromFormatted(formatted)
	if err != nil {
		panic(err)
	}
	store, err := redis.NewStoreWithOptions(rds, limiter.StoreOptions{
		Prefix:          "gull:limiter",
		CleanUpInterval: limiter.DefaultCleanUpInterval,
	})
	if err != nil {
		panic(err)
	}
	return limiter.New(store, _rate)
}
