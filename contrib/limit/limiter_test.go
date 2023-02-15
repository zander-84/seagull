package limit

import (
	"context"
	"fmt"
	goredis2 "github.com/zander-84/seagull/drive/goredis"
	"testing"
	"time"
)

func TestNewLimitContainer(t *testing.T) {
	r := goredis2.NewRdb(goredis2.Conf{
		Addr:        "172.16.86.160:6379",
		Password:    "123456",
		Db:          1,
		PoolSize:    500,
		MinIdle:     5,
		PoolTimeout: 300,
	})
	r.Start()

	l, _ := NewLimitContainer(r.Engine())
	fmt.Println(l.Get(context.Background(), "test-limiter"))
	fmt.Println(l.Allow(context.Background(), "test-limiter", 1, 0*time.Minute, 0, 2))
}

// BenchmarkNewLimitContainer-16    	   16153	     73267 ns/op
func BenchmarkNewLimitContainer(b *testing.B) {
	r := goredis2.NewRdb(goredis2.Conf{
		Addr:        "172.16.86.160:6379",
		Password:    "123456",
		Db:          1,
		PoolSize:    500,
		MinIdle:     5,
		PoolTimeout: 300,
	})
	r.Start()

	//l, _ := NewLimitContainer(r.Engine())
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.GetInt64(context.Background(), "fff")
			//l.Incr(context.Background(), "test-limiter2", 1, 10*time.Minute, 0, 100000)
		}
	})
}

func BenchmarkNewLimitContainer2(b *testing.B) {
	r := goredis2.NewRdb(goredis2.Conf{
		Addr:        "172.16.86.160:6379",
		Password:    "123456",
		Db:          1,
		PoolSize:    100,
		MinIdle:     5,
		PoolTimeout: 300,
	})
	r.Start()

	l, _ := NewLimitContainer(r.Engine())
	for i := 0; i < b.N; i++ {
		//r.GetInt64(context.Background(), "fff")
		//
		l.Allow(context.Background(), "test-limiter", 1, 10*time.Minute, 0, 10000)
	}
}
