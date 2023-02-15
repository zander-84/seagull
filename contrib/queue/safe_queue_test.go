package queue

import (
	"context"
	"errors"
	"fmt"
	"github.com/zander-84/seagull/contract"
	"github.com/zander-84/seagull/contrib/checkcode"
	"github.com/zander-84/seagull/contrib/storage"
	"github.com/zander-84/seagull/contrib/unique"
	"github.com/zander-84/seagull/drive/mongo"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestNewSafeProducer(t *testing.T) {
	mdb := mongo.NewMongo(mongo.Conf{
		Host:            "172.16.86.160",
		Port:            "27017",
		MaxPoolSize:     100,
		MinPoolSize:     10,
		MaxConnIdleTime: 5,
		Database:        "test",
		User:            "zander",
		Pwd:             "zander",
	})
	if err := mdb.Start(); err != nil {
		t.Fatal(err.Error())
	}

	m := storage.NewMongo(mdb.DB(), "queue_student_safe", 100)
	fp, err := NewSafeProducer("172.16.86.160:4150", "queue_student_safe", m, unique.New("LF", "", "aa", checkcode.NewAlpha(3)))
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

// go192 test  -timeout 1h -v  -run TestNewSafeConsumer
func TestNewSafeConsumer(t *testing.T) {
	mdb := mongo.NewMongo(mongo.Conf{
		Host:            "172.16.86.160",
		Port:            "27017",
		MaxPoolSize:     100,
		MinPoolSize:     10,
		MaxConnIdleTime: 5,
		Database:        "test",
		User:            "zander",
		Pwd:             "zander",
	})
	if err := mdb.Start(); err != nil {
		t.Fatal(err.Error())
	}
	m := storage.NewMongo(mdb.DB(), "queue_student_safe", 100)
	fc, err := NewSafeConsumer("172.16.86.160:4150", "queue_student_safe", "test", m, 3)
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
				if strings.HasPrefix(data, "1") {
					fmt.Println(data, "消费失败 ", time.Now().Format("2006-01-02 15:04:05"))
					return errors.New("data has prefix 1")
				}
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

func TestNewSafeManager(t *testing.T) {
	mdb := mongo.NewMongo(mongo.Conf{
		Host:            "172.16.86.160",
		Port:            "27017",
		MaxPoolSize:     100,
		MinPoolSize:     10,
		MaxConnIdleTime: 5,
		Database:        "test",
		User:            "zander",
		Pwd:             "zander",
	})
	if err := mdb.Start(); err != nil {
		t.Fatal(err.Error())
	}

	m := storage.NewMongo(mdb.DB(), "queue_student_safe", 100)
	fp, err := NewSafeProducer("172.16.86.160:4150", "queue_student_safe", m, unique.New("LF", "", "aa", checkcode.NewAlpha(3)))
	if err != nil {
		t.Fatal("nsq start err:" + err.Error())
	}

	_manager := NewManager(fp, "queue_student_safe", m)
	if err := _manager.SendFromBackup("LFISW148000002230116171449"); err != nil {
		t.Fatal(err.Error())
	}

	fmt.Println("success")
}
