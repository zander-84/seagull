package threadpool

import (
	"github.com/zander-84/seagull/contract"
	"sync"
)

type job struct {
	err     error
	wait    *sync.WaitGroup
	handler contract.Job
}

func newJob(Job contract.Job, wait *sync.WaitGroup) *job {
	return &job{
		err:     nil,
		wait:    wait,
		handler: Job,
	}
}
