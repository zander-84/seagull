package queue

import (
	"context"
	"errors"
	"fmt"
	"github.com/zander-84/seagull/contract"
	"github.com/zander-84/seagull/contrib/checkcode"
	"github.com/zander-84/seagull/contrib/unique"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func TestNewFastProducer(t *testing.T) {
	fp, err := NewFastProducer("172.16.86.160:4150", "test", unique.New("LF", "", "aa", checkcode.NewAlpha(3)))
	if err != nil {
		t.Fatal("nsq start err:" + err.Error())
	}

	for i := 0; i < 100; i++ {
		message := &contract.QMessage{
			Kind:       "k1",
			ForeignKey: "456",
			Data:       fmt.Sprintf("%d %s", i, time.Now().Format("2006-01-02 15:04:05")),
		}

		if err := fp.Send(message); err != nil {
			t.Fatal("send err:" + err.Error())
		}
	}

}

// go192 test  -timeout 1h -v  -run TestNewFastConsumer
func TestNewFastConsumer(t *testing.T) {
	fc, err := NewFastConsumer("172.16.86.160:4150", "test", "test", 3, nil)
	if err != nil {
		t.Fatal("nsq start err:" + err.Error())
	}
	ch := make(chan struct{}, 0)
	if err := fc.Consume(map[string]func(data string) error{
		"k1": func(data string) error {
			select {
			case <-ch:
				fmt.Println(data, "消费时间失败 ", time.Now().Format("2006-01-02 15:04:05"))
				return errors.New("exit")
			default:
				time.Sleep(time.Second)
				fmt.Println(data, "消费时间 ", time.Now().Format("2006-01-02 15:04:05"))
				return nil

			}

		},
	}); err != nil {
		t.Fatal(err.Error())
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT}...)
	select {
	case <-c:
		if err := fc.Close(context.Background()); err != nil {
			t.Fatal(err.Error())
		}
		close(ch)
	}

	fmt.Println("success")
}
