package middleware

import (
	"context"
	"github.com/zander-84/seagull/endpoint"
	"github.com/zander-84/seagull/endpoint/wraptransporter"
	"github.com/zander-84/seagull/transport"
	"github.com/zander-84/seagull/transport/http"
	"strconv"
)

func Assign(token string, pageCode string, PageSize string) endpoint.Middleware {
	return func(next endpoint.HandlerFunc) endpoint.HandlerFunc {
		return func(ctx context.Context, request interface{}) (out interface{}, err error) {
			transporter := transport.GetTransporter(ctx)
			if pageCode != "" {
				if transporter.Kind().IsHttp() {
					httpCtx := ctx.(http.Context)
					wraptransporter.SetPage(transporter, shouldStoI(httpCtx.Query().Get(pageCode)))
					wraptransporter.SetPageSize(transporter, shouldStoI(httpCtx.Query().Get(PageSize)))
				} else if transporter.Kind().IsGrpc() {
					wraptransporter.SetPage(transporter, shouldStoI(transporter.RequestHeader().Get(pageCode)))
					wraptransporter.SetPageSize(transporter, shouldStoI(transporter.RequestHeader().Get(PageSize)))
				}
			}
			if token != "" {
				wraptransporter.SetToken(transporter, transporter.RequestHeader().Get(token))
			}
			return next(ctx, request)
		}

	}
}
func shouldStoI(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return n
}
