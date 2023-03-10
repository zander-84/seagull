package cache

import (
	"context"
	"github.com/zander-84/seagull/contract/def"
	"github.com/zander-84/seagull/contrib/worker"
	memory3 "github.com/zander-84/seagull/drive/memory"
	"testing"
	"time"
)

type student struct {
	Age  int
	Name string
}

func TestNewMemoryCache(t *testing.T) {
	mem := memory3.NewMemory(memory3.Conf{
		Expiration:      10,
		CleanupInterval: 10,
	})
	mem.Start()

	c := NewMemoryCache(mem.Engine(), worker.NewProcessor())
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
	s4 := make([]student, 0)

	if err := c.BatchGetOrSet(ctx, key3, &s4, time.Hour, func(missIds []def.K) (value map[string]any, err error) {
		value = make(map[string]any, 0)
		for _, v := range missIds {
			value[v.Key] = &student{
				Age:  19,
				Name: "zander",
			}
		}
		return value, err
	}); err != nil {
		t.Fatal(err.Error())
	}
	t.Log(s4)

}
