package queue

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/nsqio/go-nsq"
	"github.com/zander-84/seagull/contract"
	"time"
)

type fastProducer struct {
	nsq    *nsq.Producer
	addr   string
	topic  string
	unique contract.Unique
}

func NewFastProducer(addr string, topic string, unique contract.Unique) (contract.QProducer, error) {
	producer, err := nsq.NewProducer(addr, nsq.NewConfig())
	if err != nil {
		return nil, err
	}
	if err := producer.Ping(); err != nil {
		return nil, err
	}
	out := new(fastProducer)
	out.nsq = producer
	out.addr = addr
	out.topic = topic
	out.unique = unique

	return out, nil
}

func (p *fastProducer) SendOrigin(message *contract.QMessage) error {
	return errors.New("the fast model does not persist the original data")
}

func (p *fastProducer) SendFromBackup(message *contract.QMessage) error {
	return errors.New("the fast model does not support SendFromBackup")
}

func (p *fastProducer) Send(message *contract.QMessage) error {
	message.SetTopic(p.topic)
	message.SetUID(p.unique.ID())
	message.SetCode(contract.QCodeWaiting)
	message.SetComeFrom(contract.QDataFromFirstTime)

	data, err := json.Marshal(message)

	if err != nil {
		return err
	}

	err = p.nsq.Publish(p.topic, data)
	if err != nil {
		err = p.nsq.Publish(p.topic, data)
		if err != nil {
			time.Sleep(time.Second)
			err = p.nsq.Publish(p.topic, data)
		}
	}
	return err
}

func (p *fastProducer) Close(ctx context.Context) error {
	p.nsq.Stop()
	return nil
}
