package cache

import (
	"context"
	"github.com/zander-84/seagull/contract/def"
	"github.com/zander-84/seagull/contrib/codec"
	"github.com/zander-84/seagull/contrib/worker"
	goredis2 "github.com/zander-84/seagull/drive/goredis"
	"testing"
	"time"
)

func TestNewRedisCache(t *testing.T) {
	r := goredis2.NewRdb(goredis2.Conf{
		Addr:        "172.16.86.160:6379",
		Password:    "123456",
		Db:          1,
		PoolSize:    500,
		MinIdle:     5,
		PoolTimeout: 300,
	})
	r.Start()

	c := NewRedisCache(r.Engine(), codec.GetCodec(codec.Json), worker.NewProcessor(), 0)
	ctx := context.Background()
	key := def.K{Key: "key"}
	s1 := student{
		Age:  18,
		Name: "zander",
	}
	if err := c.Set(ctx, key, s1, 0); err != nil {
		t.Fatal(err.Error())
	}

	s2 := new(student)
	if err := c.Get(ctx, key, s2); err != nil {
		t.Fatal(err.Error())
	}
	t.Log(s2)

	key2 := def.K{Key: "key2"}
	s3 := new(student)
	if err := c.GetOrSet(ctx, key2, s3, time.Hour, func(key def.K) (value any, err error) {
		return student{
			Age:  19,
			Name: "zander",
		}, err
	}); err != nil {
		t.Fatal(err.Error())
	}
	t.Log(s3)

	key3 := make([]def.K, 0)
	key3 = append(key3, key, key2, def.K{Key: "key3"})
	//s4 := make([]student, 0)

	//if err := c.BatchGetOrSet(ctx, key3, &s4, time.Hour, func(missIds []contract.CacheKey) (value map[contract.CacheKey]any, err error) {
	//	value = make(map[contract.CacheKey]any, 0)
	//	for _, v := range missIds {
	//		value[v] = &student{
	//			Age:  19,
	//			Name: "zander",
	//		}
	//	}
	//	fmt.Println(value)
	//	return value, err
	//}); err != nil {
	//	t.Fatal(err.Error())
	//}
	//t.Log(s4)

	if err := c.Ping(ctx); err != nil {
		t.Fatal(err.Error())
	}
}
