package threadpool

import (
	"context"
	"fmt"
	"github.com/zander-84/seagull/contract"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// go test -v  -run TestDispatcher
func TestDispatcher(t *testing.T) {
	d := NewThreadPool(Conf{
		MaxWorkers: 100,
		MaxQueues:  1000000,
	})

	if err := d.Start(); err != nil {
		t.Fatal("start Dispatcher err: ", err.Error())
	}
	wait := sync.WaitGroup{}

	var a int64
	cnt := 1000
	for i := 0; i < cnt; i++ {
		if err := d.AddJob(contract.JobFunc(func() error {
			atomic.AddInt64(&a, 1)
			time.Sleep(1 * time.Second)
			return nil
		}), &wait); err != nil {
			fmt.Println(err.Error())
		}
	}
	fmt.Println("Start Wait", time.Now().Format("2006-01-02 15:04:05"))
	wait.Wait()
	fmt.Println("Fin Wait", time.Now().Format("2006-01-02 15:04:05"))

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1000*time.Millisecond*3)
	defer cancel()
	err := d.Stop(ctx)
	if err != nil {
		t.Log("stop failed:" + err.Error())
	}
	if a == int64(cnt) {
		t.Log("success")
	} else {
		fmt.Println("a: ", a)
		fmt.Println("cnt: ", cnt)
		t.Fatal("error")
	}

}
