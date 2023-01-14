package threadpool

import (
	"errors"
	"log"
	"runtime"
)

type worker struct {
	jobChannel chan *job
	workerPool chan chan *job
	quit       chan bool
}

func newWorker(workerPool chan chan *job) *worker {
	return &worker{
		jobChannel: make(chan *job),
		workerPool: workerPool,
		quit:       make(chan bool),
	}
}

func (w *worker) start(i int) {
	go func() {
		for {
			w.workerPool <- w.jobChannel
			select {
			case _job := <-w.jobChannel:
				if _job != nil {
					func() {
						defer func() {
							if _job.wait != nil {
								_job.wait.Done()
							}
						}()
						defer func() {
							if rErr := recover(); rErr != nil {
								buf := make([]byte, 64<<10)
								n := runtime.Stack(buf, false)
								buf = buf[:n]
								log.Printf("Printf err: %v \n", rErr)
								log.Println(string(buf))
								_job.err = errors.New(string(buf))
							}
						}()
						_job.err = _job.handler.Run()
					}()
				}
			case <-w.quit:
				return
			}
		}
	}()

}

func (w *worker) stop() {
	w.quit <- true
}
