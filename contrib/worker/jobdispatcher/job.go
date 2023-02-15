package jobdispatcher

import (
	"context"
	"github.com/zander-84/seagull/think"
	"sync"
	"time"
)

type Jobs struct {
	dataSlice []*Job
	dataMap   map[string]*Job
	once      sync.Once
}

func (js *Jobs) GetSlice() []*Job {
	return js.dataSlice
}

func (js *Jobs) GetMap() map[string]*Job {
	js.once.Do(func() {
		js.dataMap = make(map[string]*Job)
		for _, v := range js.dataSlice {
			js.dataMap[v.title] = v
		}
	})
	return js.dataMap
}

func (js *Jobs) GetByTitle(title string) (*Job, error) {
	d := js.GetMap()
	e, ok := d[title]
	if ok {
		return e, nil
	} else {
		return nil, think.RecordNotFound
	}
}

type Job struct {
	startAt time.Time
	finAt   time.Time
	ctx     context.Context
	title   string
	handler func(in interface{}) (interface{}, error)
	input   interface{}
	output  interface{}
	error   error
}

func newJob(ctx context.Context, title string, input interface{}, handler func(in interface{}) (interface{}, error)) *Job {
	j := new(Job)
	j.ctx = ctx
	j.title = title
	j.handler = handler
	j.input = input
	j.error = think.UnImpl
	return j
}
func (j *Job) Title() string {
	return j.title
}
func (j *Job) Latency() time.Duration {
	return j.finAt.Sub(j.startAt)
}
func (j *Job) Result() (interface{}, error) {
	return j.output, j.error
}
