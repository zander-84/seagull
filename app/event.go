package app

import (
	"log"
)

type Event struct {
	name    string
	handler func() error
}

func NewEvent(name string, handler func() error) Event {
	return Event{
		name: name,
		handler: func() error {
			log.Println("【" + name + "】" + " run")
			err := handler()
			log.Println("【" + name + "】" + " fin")
			if err != nil {
				log.Println("【" + name + "】" + " Err: " + err.Error())
			}
			return err
		},
	}
}
