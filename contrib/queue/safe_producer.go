package queue

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/nsqio/go-nsq"
	"github.com/zander-84/seagull/contract"
	"time"
)

type safeProducer struct {
	topic  string
	addr   string
	nsq    *nsq.Producer
	mongo  contract.Mongo
	unique contract.Unique
}

func NewSafeProducer(addr string, topic string, mongo contract.Mongo, unique contract.Unique) (contract.QProducer, error) {
	producer, err := nsq.NewProducer(addr, nsq.NewConfig())
	if err != nil {
		return nil, err
	}
	if err := producer.Ping(); err != nil {
		return nil, err
	}

	out := new(safeProducer)
	out.nsq = producer
	out.addr = addr
	out.mongo = mongo
	out.topic = topic
	out.unique = unique
	return out, nil
}

func (p *safeProducer) SendOrigin(message *contract.QMessage) error {
	message.SetCode(contract.QCodeOrigin)
	_, err := p.mongo.Create(message)
	return err
}

func (p *safeProducer) Send(message *contract.QMessage) error {
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

	var dbErr error
	if message.CreatedAt > 0 {
		dbErr = p.mongo.ReplaceOneByKv("uid", message.UID, 0, message)
		if dbErr != nil {
			dbErr = p.mongo.ReplaceOneByKv("uid", message.UID, 0, message)
		}
	} else {
		_, dbErr = p.mongo.Create(message)
		if dbErr != nil {
			_, dbErr = p.mongo.Create(message)
		}
	}

	err = p.nsq.Publish(p.topic, data)
	if err != nil {
		err = p.nsq.Publish(p.topic, data)
		if err != nil {
			time.Sleep(time.Second)
			err = p.nsq.Publish(p.topic, data)
		}
	}
	if err == nil || dbErr == nil {
		return nil
	}
	if err != nil && dbErr != nil {
		return err
	}
	return err
}

func (p *safeProducer) SendFromBackup(message *contract.QMessage) error {
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

func (p *safeProducer) Close(ctx context.Context) error {
	p.nsq.Stop()
	return nil
}
