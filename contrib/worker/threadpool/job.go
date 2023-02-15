package threadpool

import (
	"github.com/zander-84/seagull/contract/def"
	"sync"
)

type job struct {
	err     error
	wait    *sync.WaitGroup
	handler def.Job
}

func newJob(Job def.Job, wait *sync.WaitGroup) *job {
	return &job{
		err:     nil,
		wait:    wait,
		handler: Job,
	}
}
