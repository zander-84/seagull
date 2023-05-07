package project

var RmcTpl = `package transport

import (
	"context"
	"github.com/zander-84/seagull/endpoint"
	"github.com/zander-84/seagull/endpoint/middleware"
	"github.com/zander-84/seagull/endpoint/middleware/cors"
	"github.com/zander-84/seagull/think"
	"github.com/zander-84/seagull/transport/grpc"
	"github.com/zander-84/seagull/transport/http"
	"${project}/apps/${server}/internal/endpoint/hello"
	"${project}/apps/${server}/internal/pkg"
)

func NewRmc(mode think.Mode) endpoint.Rmc {
	resource := endpoint.NewRmc()
	resource = resource.Use(endpoint.OptErrorEncoder(endpoint.WrapError(map[endpoint.Kind]func(ctx context.Context, err error) error{
		endpoint.Http: func(ctx context.Context, err error) error {
			err = ctx.(http.Context).ErrorEncoder(err, mode.IsProd())
			return err
		},
		endpoint.Grpc: func(ctx context.Context, err error) error {
			err = ctx.(grpc.Context).ErrorEncoder(err, mode.IsProd())
			return err
		},
	}))).Use(endpoint.OptMW(middleware.Recover())).Use(endpoint.OptMW(cors.New(pkg.GinCors(mode)))).
		Use(endpoint.OptMW(middleware.Assign("Token", "Page", "PageSize")))


	resource.Endpoint(endpoint.MethodGet, "/", hello.Hello, endpoint.Codecs{
		endpoint.Http: {Dec: hello.HelloCodec{}.HttpGetDecode, Enc: hello.HelloCodec{}.HttpGetEncode},
		endpoint.Grpc: {Dec: hello.HelloCodec{}.GrpcDecode, Enc: hello.HelloCodec{}.GrpcEncode},
	})

	return resource
}

`
