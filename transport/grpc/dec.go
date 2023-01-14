package grpc

import (
	"github.com/zander-84/seagull/think"
)

func Dec(request any, out any) error {
	in, ok := request.(func(interface{}) error)
	if !ok {
		return think.ErrAlert("something err : request error")
	}
	if err := in(out); err != nil {
		return think.ErrAlert(err.Error())
	}
	return nil
}
