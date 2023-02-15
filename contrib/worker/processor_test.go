package worker

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestNewGo(t *testing.T) {
	g := NewProcessor()

	//if err := g.RunWithTimeout(func() error {
	//	time.Sleep(time.Second)
	//	return nil
	//}, time.Second); err != nil {
	//	t.Fatal(err.Error())
	//}
	//
	//ctx1, cancel1 := context.WithTimeout(context.Background(), time.Second*10)
	//go func() {
	//	time.Sleep(time.Second)
	//	//fmt.Println("do cancel")
	//	//cancel1()
	//}()
	//defer cancel1()
	//ctx2, cancel2 := context.WithTimeout(context.Background(), time.Second*10)
	//ctx3, cancel3 := context.WithTimeout(context.Background(), time.Second*2)
	//defer cancel2()
	//defer cancel3()
	//
	//if data, err := g.GoListenCtx([]context.Context{ctx1, ctx2, ctx3}, func() (any, error) {
	//	fmt.Println("GoListenCtx")
	//	time.Sleep(time.Second * 3)
	//	fmt.Println("GoListenCtx2")
	//
	//	return 123, nil
	//}); err != nil {
	//	t.Fatal(err.Error())
	//} else {
	//	fmt.Println(data)
	//}

	w := sync.WaitGroup{}
	g.Go(func() {
		fmt.Println("before sleep")
		time.Sleep(time.Second * 10)
		fmt.Println("after sleep")
	}, &w)
	w.Wait()
	fmt.Println("success")
	//if err := g.Wait(time.Second * 2); err != nil {
	//	t.Fatal(err.Error())
	//} else {
	//	fmt.Println("success")
	//}
}
