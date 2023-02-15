package queue

import (
	"context"
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"github.com/zander-84/seagull/contract"
	"sync"
)

type fastConsumer struct {
	topic        string
	nsq          *nsq.Consumer
	consumerAddr string
	recordFailed func(message *contract.QMessage)
	maxWorker    int
	waiter       sync.WaitGroup
}

func NewFastConsumer(addr string, topic string, ch string, maxWorker int, recordFailed func(message *contract.QMessage)) (contract.QConsumer, error) {
	conf := nsq.NewConfig()
	conf.MaxInFlight = maxWorker
	consumer, err := nsq.NewConsumer(topic, ch, conf)
	if err != nil {
		return nil, err
	}
	out := new(fastConsumer)
	out.nsq = consumer
	out.maxWorker = maxWorker
	out.consumerAddr = addr
	out.recordFailed = recordFailed
	out.topic = topic
	return out, nil
}

func (f *fastConsumer) Consume(workers map[string]func(data string) error) error {
	f.nsq.AddConcurrentHandlers(nsq.HandlerFunc(func(message *nsq.Message) error {
		f.waiter.Add(1)
		defer f.waiter.Done()
		body := message.Body
		qMess := new(contract.QMessage)
		if err := json.Unmarshal(body, qMess); err != nil {
			return nil
		}
		worker, ok := workers[qMess.Kind]
		if !ok {
			// 兜底worker
			worker, ok = workers[""]
		}
		if !ok {
			return nil
		}
		err := worker(qMess.Data)
		if err == nil {
			message.Finish()
		}

		if err != nil && f.recordFailed != nil {
			f.recordFailed(qMess)
		}
		return err
	}), f.maxWorker)
	return f.nsq.ConnectToNSQD(f.consumerAddr)
}

func (f *fastConsumer) Close(ctx context.Context) error {
	fin := make(chan struct{}, 1)
	go func() {
		f.nsq.Stop()
		f.waiter.Wait()
		fin <- struct{}{}
	}()
	done := ctx.Done()
	if done != nil {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-fin:
			return nil
		}
	} else {
		select {
		case <-fin:
			return nil
		}
	}
}
