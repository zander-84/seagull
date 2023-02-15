package crontab

import (
	"fmt"
	"github.com/zander-84/seagull/contract"
	"github.com/zander-84/seagull/contract/def"
	"testing"
	"time"
)

// go test -v  -run TestRobfig
func TestRobfig(t *testing.T) {
	c := NewCrontab(Conf{})
	if err := c.Start(); err != nil {
		t.Fatal("start cron err: ", err.Error())
	}
	testAdd(t, c)
	testRemove(t, c)
	testAdd(t, c)

	if err := c.StartJobs(); err != nil {
		t.Fatal("StartJobs error ", err.Error())
	}

	//if err := c.StartJobs(); err != nil {
	//	t.Fatal("StartJobs error ", err.Error())
	//}

	time.Sleep(2 * time.Minute)

	t.Log("success")
}

func testAdd(t *testing.T, c *Crontab) {
	if err := c.AddJob(contract.CronJob{
		ID:   "test1",
		Desc: "测试1",
		Spec: "* * * * * *",
		Cmd: def.JobFunc(func() error {
			fmt.Println("hello world")
			return nil
		}),
	}); err != nil {
		t.Fatal("add error ", err.Error())
	}

	if err := c.AddJob(contract.CronJob{
		ID:   "test2",
		Desc: "测试2",
		Spec: "* * * * * *",
		Cmd: def.JobFunc(func() error {
			fmt.Println("hello world 2222")
			return nil
		}),
	}); err != nil {
		t.Fatal("add error ", err.Error())
	}
}
func testRemove(t *testing.T, c *Crontab) {
	if err := c.RemoveJob("test1"); err != nil {
		t.Fatal("add error ", err.Error())
	}
	if err := c.RemoveJob("test2"); err != nil {
		t.Fatal("add error ", err.Error())
	}
}
