package threadpool

import (
	"context"
	"errors"
	"github.com/zander-84/seagull/contract/def"
	"github.com/zander-84/seagull/think"
	"sync"
	"sync/atomic"
)

type Dispatcher struct {
	worker     []*worker      // 工人
	workerPool chan chan *job // 真正任务派发队列
	jobChannel chan *job
	conf       Conf
	err        error
	lock       sync.Mutex
	once       int64
	exit       chan struct{}
}

func NewThreadPool(conf Conf) *Dispatcher {
	var this = new(Dispatcher)
	this.init(conf)
	return this
}

func (d *Dispatcher) init(conf Conf) {
	d.conf = conf.SetDefault()
	atomic.StoreInt64(&d.once, 0)
	d.err = think.UnImpl
}

func (d *Dispatcher) Start() error {
	d.lock.Lock()
	defer d.lock.Unlock()

	if atomic.CompareAndSwapInt64(&d.once, 0, 1) {
		d.worker = make([]*worker, d.conf.MaxWorkers)
		d.workerPool = make(chan chan *job, d.conf.MaxWorkers)
		d.jobChannel = make(chan *job, d.conf.MaxQueues)
		d.exit = make(chan struct{}, 1)
		d.run()
		d.err = nil
	}

	return d.err
}

func (d *Dispatcher) Stop(ctx context.Context) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	done := ctx.Done()
	if done == nil {
		d.stop()
		return nil
	}
	fin := make(chan struct{}, 1)
	go func() {
		d.stop()
		fin <- struct{}{}
	}()
	select {
	case <-done:
		return ctx.Err()
	case <-fin:
		return nil
	}
}
func (d *Dispatcher) stop() {
	if atomic.CompareAndSwapInt64(&d.once, 1, 2) {
		close(d.jobChannel)
		<-d.exit

		for _, w := range d.worker {
			w.stop()
		}
	}
}

func (d *Dispatcher) run() {
	for i := 0; i < len(d.worker); i++ {
		d.worker[i] = newWorker(d.workerPool)
		d.worker[i].start(i)
	}
	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for _job := range d.jobChannel {
		jobChannel := <-d.workerPool
		jobChannel <- _job
	}
	d.exit <- struct{}{}
}

func (d *Dispatcher) AddJobUnsafe(job def.Job, wait *sync.WaitGroup) (err error) {
	if d.conf.MaxQueues <= len(d.jobChannel) {
		return errors.New("system busyness")
	}
	return d.addJobWait(job, wait)
}

func (d *Dispatcher) AddJob(job def.Job, wait *sync.WaitGroup) (err error) {
	return d.addJobWait(job, wait)
}

func (d *Dispatcher) addJobWait(job def.Job, wait *sync.WaitGroup) (err error) {
	if wait != nil {
		wait.Add(1)
	}
	err = d.addJob(job, wait)
	if err != nil {
		if wait != nil {
			wait.Done()
		}
	}
	return err
}

func (d *Dispatcher) addJob(job def.Job, wait *sync.WaitGroup) (err error) {
	defer func() {
		if recoverErr := recover(); recoverErr != nil {
			err = errors.New("queue already exited")
		}
	}()

	d.jobChannel <- newJob(job, wait)
	return nil

}
