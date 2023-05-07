package grpc

import (
	"context"
	"github.com/zander-84/seagull/think"
	"github.com/zander-84/seagull/tool/conv"
	"github.com/zander-84/seagull/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"
)

var _ Context = (*wrapper)(nil)

// Context is an HTTP Context.
type Context interface {
	context.Context
	ErrorEncoder(err error, isProdEnv bool) error
	Encoder(v any) func() (any, error)
	//RecoverErr(isProdEnv bool) (context.Context, error)
}

func NewGrpcContext(ctx context.Context, transporter transport.Transporter) Context {
	w := &wrapper{ctx: ctx, transporter: transporter}
	return w
}

type wrapper struct {
	ctx         context.Context
	transporter transport.Transporter
}

func (c *wrapper) Deadline() (time.Time, bool) {
	return c.ctx.Deadline()
}

func (c *wrapper) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *wrapper) Err() error {
	return c.ctx.Err()
}

func (c *wrapper) Value(key interface{}) interface{} {
	return c.ctx.Value(key)
}

func (c *wrapper) ErrorEncoder(err error, isProdEnv bool) error {
	thinkErr := think.FromError(err)
	var data any
	if isProdEnv && think.IsErrSystemSpace(thinkErr) {
		data = thinkErr.Response.Data
		thinkErr.Response.Data = thinkErr.Code.ToString()
	}
	err = thinkErr.ToGrpcErr()

	if isProdEnv && think.IsErrSystemSpace(thinkErr) {
		thinkErr.Response.Data = data
	}

	return err
}

func (c *wrapper) Encoder(v any) func() (any, error) {
	return func() (any, error) {
		md := map[string]string{}

		c.transporter.ReplyHeader().Foreach(func(k, v string) error {
			md[k] = v
			return nil
		})

		md["code"] = conv.IntegerToStr(int32(c.transporter.Code()))
		md["bizCode"] = c.transporter.BizCode()

		err := grpc.SetHeader(c.ctx, metadata.New(md))
		if err != nil {
			return nil, think.ErrSystemSpace(err.Error())
		}

		return v, nil
	}

}
