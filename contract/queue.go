package contract

import "context"

type QMessage struct {
	ID     string `json:"id"`
	Kind   string `json:"kind"`
	Status int    `json:"status"`
	Data   string `json:"data"`
}

type QProducer interface {
	Send(topic string, message *QMessage) error
	Close(ctx context.Context) error
}

type QConsumer interface {
	Fetch(topic string, message *QMessage) error
	Close(ctx context.Context) error
}
