package zap

import (
	"context"
	"fmt"
	"github.com/zander-84/seagull/contract/def"
	"io"
	"log"
	"testing"
)

func TestNewZapLog(t *testing.T) {
	log.Println("abc")

	l, cancel, err := NewZapLog(Conf{
		Level:     "info",
		Name:      "test",
		AddCaller: false,
		ConsoleHook: struct {
			Enable bool
		}{
			Enable: true,
		},
		FileHook: struct {
			Enable     bool
			Path       string
			MaxAge     int
			MaxBackups int
			MaxSize    int
		}{
			Enable:     true,
			Path:       "./logs",
			MaxAge:     10,
			MaxBackups: 10,
			MaxSize:    50,
		},
	}, []io.Writer{Writer(func(p []byte) (n int, err error) {
		fmt.Println("Writer: ", string(p))
		return 0, nil
	})})
	if err != nil {
		t.Fatal(err.Error())
	}
	defer cancel(context.Background())
	l.Debug(def.E{Key: "content", Value: "Debug"}, def.E{Key: "id", Value: 123})
	l.Info(def.E{Key: "content", Value: "Info"}, def.E{Key: "id", Value: 123})
	l.Error(def.E{Key: "msg", Value: "Error"}, def.E{Key: "id", Value: 123})
	log.Println("abc")
}
