package middleware

import (
	"context"
	"github.com/zander-84/seagull/endpoint"
	"github.com/zander-84/seagull/think"
	"github.com/zander-84/seagull/tool/conv"
	"github.com/zander-84/seagull/transport/http"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func BeforeResponse() endpoint.Middleware {
	return func(next endpoint.HandlerFunc) endpoint.HandlerFunc {
		return func(ctx context.Context, request interface{}) (out interface{}, err error) {
			out, err = next(ctx, request)

			transporter := endpoint.GetTransporter(ctx)

			if transporter.Kind().IsHttp() {
				httpCtx := ctx.(http.Context)

				transporter.ReplyHeader().Foreach(func(k, v string) error {
					httpCtx.Response().Header().Set(k, v)
					return nil
				})

			} else if transporter.Kind().IsGrpc() {
				md := map[string]string{}

				transporter.ReplyHeader().Foreach(func(k, v string) error {
					md[k] = v
					return nil
				})

				md["code"] = conv.IntegerToStr(int32(transporter.Code()))
				md["bizCode"] = transporter.BizCode()

				err = grpc.SetHeader(ctx, metadata.New(md))
				if err != nil {
					err = think.ErrSystemSpace(err.Error())
				}
			}

			return out, err
		}
	}
}
