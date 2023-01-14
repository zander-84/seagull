package crontab

import (
	"github.com/robfig/cron/v3"
	"github.com/zander-84/seagull/contract"
	"log"
	"runtime"
)

type job struct {
	id      cron.EntryID
	cronJob contract.CronJob
}

func newJob(_job contract.CronJob) *job {
	out := new(job)
	out.cronJob = _job
	return out
}

func (j *job) Run() {
	defer func() {
		if rErr := recover(); rErr != nil {
			buf := make([]byte, 64<<10)
			n := runtime.Stack(buf, false)
			buf = buf[:n]
			log.Printf("Printf err: %v \n", rErr)
			log.Println(string(buf))
		}
	}()
	_ = j.cronJob.Cmd.Run()
}
