package http_gin

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/zander-84/seagull/endpoint"
	"github.com/zander-84/seagull/transport"
	"github.com/zander-84/seagull/transport/http"
	http2 "net/http"
)

var _ http2.Handler = (*Router)(nil)

type Router struct {
	ginEngine *gin.Engine
}

func NewRouter(ginEngine *gin.Engine) *Router {
	router := new(Router)
	router.ginEngine = ginEngine

	return router
}

func (r *Router) ServeHTTP(res http2.ResponseWriter, req *http2.Request) {
	r.ginEngine.Handler().ServeHTTP(res, req)
}
func (r *Router) Endpoint(kind transport.Kind, fullPath endpoint.Path, e endpoint.HandlerFunc) {

	switch fullPath.Method() {
	case transport.MethodGet:
		r.ginEngine.GET(fullPath.Path(), func(ctx *gin.Context) {
			_, _ = e(initCtx(ctx, kind, fullPath), nil)
		})
	case transport.MethodHead:
		r.ginEngine.OPTIONS(fullPath.Path(), func(ctx *gin.Context) {
			_, _ = e(initCtx(ctx, kind, fullPath), nil)
		})
	case transport.MethodPost:
		r.ginEngine.POST(fullPath.Path(), func(ctx *gin.Context) {
			_, _ = e(initCtx(ctx, kind, fullPath), nil)
		})
	case transport.MethodPut:
		r.ginEngine.PUT(fullPath.Path(), func(ctx *gin.Context) {
			_, _ = e(initCtx(ctx, kind, fullPath), nil)
		})
	case transport.MethodPatch:
		r.ginEngine.PATCH(fullPath.Path(), func(ctx *gin.Context) {
			_, _ = e(initCtx(ctx, kind, fullPath), nil)
		})
	case transport.MethodDelete:
		r.ginEngine.DELETE(fullPath.Path(), func(ctx *gin.Context) {
			_, _ = e(initCtx(ctx, kind, fullPath), nil)

		})
	case transport.MethodOptions:
		r.ginEngine.OPTIONS(fullPath.Path(), func(ctx *gin.Context) {
			_, _ = e(initCtx(ctx, kind, fullPath), nil)
		})
	default:

	}
}

func initCtx(ctx *gin.Context, kind transport.Kind, fullPath endpoint.Path) context.Context {
	endpointCtxVal := transport.NewTransporter(kind, transport.NewHeader(ctx.Request.Header), fullPath.FullPath(), fullPath.Method())
	ctx.Request = ctx.Request.WithContext(transport.WithContext(ctx.Request.Context(), endpointCtxVal))
	return http.NewHttpContext(ctx.Writer, ctx.Request, &proxy{ctx: ctx}, endpointCtxVal)
}
