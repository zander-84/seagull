package queue

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/nsqio/go-nsq"
	"github.com/zander-84/seagull/contract"
	"time"
)

type friendlyProducer struct {
	topic  string
	addr   string
	nsq    *nsq.Producer
	mongo  contract.Mongo
	unique contract.Unique
}

func NewFriendlyProducer(addr string, topic string, mongo contract.Mongo, unique contract.Unique) (contract.QProducer, error) {
	producer, err := nsq.NewProducer(addr, nsq.NewConfig())
	if err != nil {
		return nil, err
	}
	if err := producer.Ping(); err != nil {
		return nil, err
	}

	out := new(friendlyProducer)
	out.nsq = producer
	out.addr = addr
	out.mongo = mongo
	out.topic = topic
	out.unique = unique
	return out, nil
}

func (p *friendlyProducer) SendOrigin(message *contract.QMessage) error {
	return errors.New("the friendly model does not persist the original data")
}

func (p *friendlyProducer) Send(message *contract.QMessage) error {
	message.SetTopic(p.topic)
	message.SetUID(p.unique.ID())
	message.SetCode(contract.QCodeWaiting)
	message.SetComeFrom(contract.QDataFromFirstTime)

	if err := message.Valid(); err != nil {
		return err
	}
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	message.SetCode(contract.QCodeWaiting)
	err = p.nsq.Publish(p.topic, data)
	if err != nil {
		err = p.nsq.Publish(p.topic, data)
		if err != nil {
			time.Sleep(time.Second)
			err = p.nsq.Publish(p.topic, data)
		}
	}

	if err != nil {
		message.Code = contract.QCodeInQFail
		message.Reason = err.Error()
		message.ComeFrom = contract.QDataFromBackup
		_, err = p.mongo.Create(message)
		if err != nil {
			_, err = p.mongo.Create(message)
		}
		if err != nil {
			return err
		}
	}
	return err
}

func (p *friendlyProducer) SendFromBackup(message *contract.QMessage) error {
	if message.ComeFrom != contract.QDataFromBackup {
		return errors.New("data must come from backup")
	}
	if err := message.Valid(); err != nil {
		return err
	}
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return p.nsq.Publish(p.topic, data)
}

func (p *friendlyProducer) Close(ctx context.Context) error {
	p.nsq.Stop()
	return nil
}
