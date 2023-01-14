package lock

import (
	"context"
	"fmt"
	"github.com/zander-84/seagull/contract"
	"github.com/zander-84/seagull/contrib/checkcode"
	"github.com/zander-84/seagull/contrib/unique"
	"github.com/zander-84/seagull/contrib/worker"
	goredis2 "github.com/zander-84/seagull/drive/goredis"
	"sync"
	"testing"
	"time"
)

func TestNewStandaloneLock(t *testing.T) {
	l, cancel := NewStandaloneLock(unique.New("LF", "", "aa", checkcode.NewAlpha(3)), worker.NewProcessor(), time.Second*5)
	defer cancel()
	s1 := time.Now()
	w := &sync.WaitGroup{}
	w.Add(1)

	w.Add(1)
	go func() {
		time.Sleep(time.Second * 10)
		fmt.Println("go  begin do3")
		_do3(l, w)
		fmt.Println("go end do3")
	}()

	_do3(l, w)

	for i := 0; i < 10; i++ {
		w.Add(2)
		go _do(l, w, i)
		go _do(l, w, i)
	}
	w.Wait()
	s2 := time.Now()
	fmt.Println(s2.Sub(s1).Milliseconds())
}
func TestNewRedisLocker(t *testing.T) {
	r := goredis2.NewRdb(goredis2.Conf{
		Addr:        "172.16.86.160:6379",
		Password:    "123456",
		Db:          1,
		PoolSize:    500,
		MinIdle:     5,
		PoolTimeout: 300,
	})
	r.Start()
	l, err := NewRedisLocker(r.Engine(), unique.New("LF", "", "aa", checkcode.NewAlpha(3)), worker.NewProcessor(), time.Second*5)
	if err != nil {
		t.Fatal(err.Error())
	}

	s1 := time.Now()
	w := &sync.WaitGroup{}

	w.Add(1)
	_do3(l, w)

	for i := 0; i < 10; i++ {
		w.Add(2)
		go _do(l, w, 1)
		go _do(l, w, 1)
	}
	w.Add(1)
	_do2(l, w)
	w.Wait()
	s2 := time.Now()
	fmt.Println(s2.Sub(s1).Milliseconds())
	time.Sleep(time.Second)
}
func TestNewRedisV2Locker(t *testing.T) {
	r := goredis2.NewRdb(goredis2.Conf{
		Addr:        "172.16.86.160:6379",
		Password:    "123456",
		Db:          1,
		PoolSize:    500,
		MinIdle:     5,
		PoolTimeout: 300,
	})
	r.Start()

	l, cancel, err := NewRedisWaitLocker(r.Engine(), unique.New("LF", "", "aa", checkcode.NewAlpha(3)), worker.NewProcessor(), "testchan", time.Second*5, time.Minute)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer cancel()
	s1 := time.Now()
	w := &sync.WaitGroup{}
	//w.Add(1)
	//_do3(l, w)

	for i := 0; i < 10; i++ {
		w.Add(2)
		go _do(l, w, 1)
		go _do(l, w, 1)
	}

	w.Add(1)
	_do2(l, w)
	w.Wait()
	s2 := time.Now()
	fmt.Println(s2.Sub(s1).Milliseconds())
	time.Sleep(time.Second)

}
func _do(locker contract.Locker, w *sync.WaitGroup, num int) {
	key := "key" + fmt.Sprintf("%d", num)
	fmt.Println("key", key)
	defer w.Done()
	locked, err := locker.Lock(context.Background(), key, time.Second*20)
	if err != nil {
		if err == contract.LockFailed {
			fmt.Println("lock err")
		}
		return
	}

	defer locked.Release(context.Background())
	fmt.Println("do:", key)
	time.Sleep(time.Second)
}

func _do2(locker contract.Locker, w *sync.WaitGroup) {

	defer w.Done()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second/2)
	defer cancel()
	locked, err := locker.Lock(ctx, "key2", time.Second*20)
	if err != nil {
		if err == contract.LockFailed {
			fmt.Println("lock err")
		}
		return
	}

	defer locked.Release(context.Background())
	fmt.Println("do2")
	time.Sleep(time.Second)

}
func _do3(locker contract.Locker, w *sync.WaitGroup) {

	defer w.Done()

	locked, err := locker.Lock(context.Background(), "leaser", time.Second*20)
	if err != nil {
		if err == contract.LockFailed {
			fmt.Println("lock err")
		}
		return
	}
	defer locked.Release(context.Background())
	fmt.Println("do2")
	time.Sleep(time.Second * 30)

}
