package grpc

import (
	"context"
	"github.com/zander-84/seagull/think"
	"time"
)

var _ Context = (*wrapper)(nil)

// Context is an HTTP Context.
type Context interface {
	context.Context

	ErrorEncoder(err error, isProdEnv bool) error
	//RecoverErr(isProdEnv bool) (context.Context, error)
}

func NewGrpcContext(ctx context.Context) Context {
	w := &wrapper{ctx: ctx}
	return w
}

type wrapper struct {
	ctx context.Context
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
