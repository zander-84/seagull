package queue

import (
	"context"
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"github.com/zander-84/seagull/contract"
)

type safeConsumer struct {
	topic        string
	nsq          *nsq.Consumer
	mongo        contract.Mongo
	consumerAddr string
	maxWorker    int
}

func NewSafeConsumer(addr string, topic string, ch string, mongo contract.Mongo, maxWorker int) (contract.QConsumer, error) {
	conf := nsq.NewConfig()
	conf.MaxInFlight = maxWorker
	consumer, err := nsq.NewConsumer(topic, ch, conf)
	if err != nil {
		return nil, err
	}

	out := new(safeConsumer)
	out.nsq = consumer
	out.maxWorker = maxWorker
	out.consumerAddr = addr
	out.mongo = mongo
	out.topic = topic
	return out, nil
}

func (f *safeConsumer) Consume(workers map[string]func(data string) error) error {
	f.nsq.AddConcurrentHandlers(nsq.HandlerFunc(func(message *nsq.Message) error {
		body := message.Body
		qMess := new(contract.QMessage)
		if err := json.Unmarshal(body, qMess); err != nil {
			return nil
		}
		// 原始数据不进行任何消费
		if qMess.Code == contract.QCodeOrigin {
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

		if err := worker(qMess.Data); err != nil {
			//消费失败 第一次就记录到备份数据中  第二次就更新备份数据
			qMess.SetReason(err.Error())
			qMess.SetCode(contract.QCodeConsumeFail)
			err = f.mongo.ReplaceOneByKv("uid", qMess.UID, 0, qMess)
			if err == nil {
				message.Finish()
			}
			return err
		}

		//消费成功
		_ = f.mongo.DelOneByKv("uid", qMess.UID)
		message.Finish()
		return nil
	}), f.maxWorker)
	return f.nsq.ConnectToNSQD(f.consumerAddr)
}

func (f *safeConsumer) Close(ctx context.Context) error {
	f.nsq.Stop()
	return nil
}
