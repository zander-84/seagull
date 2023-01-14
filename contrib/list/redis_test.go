package list

import (
	"context"
	"fmt"
	"github.com/zander-84/seagull/contrib/codec"
	goredis2 "github.com/zander-84/seagull/drive/goredis"
	"testing"
)

type student struct {
	Age  int
	Name string
}

func TestNewRedisList(t *testing.T) {
	r := goredis2.NewRdb(goredis2.Conf{
		Addr:        "172.16.86.160:6379",
		Password:    "123456",
		Db:          1,
		PoolSize:    500,
		MinIdle:     5,
		PoolTimeout: 300,
	})
	r.Start()

	list := NewRedisList(r.Engine(), codec.GetCodec(codec.Json))

	ctx := context.Background()
	key := "test-list"
	//s1 := student{
	//	Age:  18,
	//	Name: "zander",
	//}
	//if err := list.LeftPush(ctx, key, s1); err != nil {
	//	t.Fatal(err.Error())
	//}

	//s3 := new(student)
	//if err := list.LeftPop(ctx, key, s3, func(data any) (bool, error) {
	//	return true, nil
	//}); err != nil {
	//	t.Fatal(err.Error())
	//}
	//fmt.Println(s3)
	//
	//s4 := new(student)
	//if err := list.RightPop(ctx, key, s4, func(data any) (bool, error) {
	//	return true, nil
	//}); err != nil {
	//	t.Fatal(err.Error())
	//}
	//fmt.Println(s4)

	//locker := lock.NewRedisLocker(r.Engine(), unique.New("a", "", "a", checkcode.NewAlpha()))

	//go list.LeftBatchPush(ctx, key, locker, 2, func(int2 int) (any, error) {
	//	return student{
	//		Age:  18,
	//		Name: "zander5",
	//	}, nil
	//})
	//
	//list.RightBatchPush(ctx, key, locker, 2, func(i int) (any, error) {
	//	return student{
	//		Age:  18,
	//		Name: "zander4",
	//	}, nil
	//})
	//for i := 0; i < 20; i++ {
	//	s2 := student{
	//		Age:  i,
	//		Name: "zander",
	//	}
	//	if err := list.RightPush(ctx, key, s2); err != nil {
	//		t.Fatal(err.Error())
	//	}
	//}

	//data := make([]student, 0)
	//if err := list.Range(ctx, key, &data, 1, 10); err != nil {
	//	t.Fatal(err)
	//} else {
	//	fmt.Println(data)
	//}
	//
	//data1 := make([]student, 0)
	//if err := list.Range(ctx, key, &data1, -10, 4); err != nil {
	//	t.Fatal(err)
	//} else {
	//	fmt.Println(data1)
	//}

	data2 := new(student)
	if err := list.RightPop(ctx, key, data2, func(data any) (bool, error) {
		return false, nil
	}); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(data2)
	}
}
