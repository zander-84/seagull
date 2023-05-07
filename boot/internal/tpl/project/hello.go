package project

var HelloTpl = `package hello

import (
	"context"
	"github.com/zander-84/seagull/endpoint"
)

func Hello(ctx context.Context, request interface{}) (response interface{}, err error) {
	in := request.(*HelloCodec)
	transporter := endpoint.GetTransporter(ctx)
	transporter.ReplyHeader().Set("Trace-Id", "123")
	return in, nil
}

`
var HelloCodecTpl = `package hello

import (
	"context"
	"github.com/zander-84/seagull/pbs"
	"github.com/zander-84/seagull/think"
	"github.com/zander-84/seagull/transport/grpc"
	"github.com/zander-84/seagull/transport/http"
)

type HelloCodec struct {
	In   *pbs.Request
	Name string
	Age  int
}

func (HelloCodec) HttpGetDecode(ctx context.Context, request interface{}) (interface{}, error) {
	return &HelloCodec{
		Name: "zander",
		Age:  18,
	}, nil
}

func (HelloCodec) HttpGetEncode(ctx context.Context, request interface{}) (any, error) {
	httpCtx := ctx.(http.Context)

	resp := think.NewSuccessResp(request)
	return httpCtx.JSON(resp.Code.HttpCode(), resp), nil
}

func (HelloCodec) GrpcDecode(ctx context.Context, request interface{}) (interface{}, error) {
	in := new(pbs.Request)
	if err := grpc.Dec(request, in); err != nil {
		return nil, err
	}

	return &HelloCodec{
		In: in,
	}, nil
}

func (HelloCodec) GrpcEncode(ctx context.Context, request interface{}) (interface{}, error) {
	grpcCtx := ctx.(grpc.Context)

	data := new(pbs.Response)
	data.AdminName = "zander"
	return grpcCtx.Encoder(data), nil
}
`
