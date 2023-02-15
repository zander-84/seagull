package queue

import (
	"context"
	"github.com/zander-84/seagull/contract"
)

type manager struct {
	topic    string
	producer contract.QProducer
	mongo    contract.Mongo
}

func NewManager(producer contract.QProducer, topic string, mongo contract.Mongo) contract.QManager {
	out := new(manager)
	out.producer = producer
	out.mongo = mongo
	out.topic = topic

	return out
}

func (f *manager) Messages(searchMeta contract.SearchMeta, searchParams contract.MongoBuilder) (message []contract.QMessage, cnt *int64, err error) {
	err = f.mongo.Search(searchMeta, searchParams, &message, cnt)
	return message, cnt, err
}

func (f *manager) AdjustMessage(message *contract.QMessage) error {
	return f.mongo.ReplaceOneByKv("uid", message.UID, message.Version, message)
}

func (f *manager) SendFromBackup(Id string) error {
	mess := new(contract.QMessage)
	if err := f.mongo.FindOneByField("uid", Id, mess); err != nil {
		return err
	}
	mess.SetCode(contract.QCodeWaiting)
	err := f.mongo.UpdatePartByKv("uid", mess.UID, 0, mess)
	if err != nil {
		return err
	}

	mess.SetComeFrom(contract.QDataFromBackup)
	if err = f.producer.SendFromBackup(mess); err != nil {
		mess.SetCode(contract.QCodeInQFail)
		mess.SetReason(err.Error())
		err = f.mongo.UpdatePartByKv("uid", mess.UID, 0, mess)
		if err != nil {
			return err
		}
	}

	return err
}

func (f *manager) ReleaseMessage(id string) error {
	return f.mongo.DelOneByKv("id", id)
}

func (f *manager) Close(ctx context.Context) error {
	return f.producer.Close(ctx)
}
