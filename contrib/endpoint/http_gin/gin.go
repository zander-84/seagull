package http_gin

import (
	"github.com/gin-gonic/gin"
	"github.com/zander-84/seagull/endpoint"
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
func (r *Router) Endpoint(kind endpoint.Kind, fullPath endpoint.Path, e endpoint.HandlerFunc) {

	switch fullPath.Method() {
	case endpoint.MethodGet:
		r.ginEngine.GET(fullPath.Path(), func(ctx *gin.Context) {
			initCtx(ctx, kind, fullPath)
			gullCtx := http.NewHttpContext(ctx.Writer, ctx.Request, &proxy{ctx: ctx})
			_, _ = e(gullCtx, nil)
			gullCtx.Push()
		})
	case endpoint.MethodHead:
		r.ginEngine.OPTIONS(fullPath.Path(), func(ctx *gin.Context) {
			initCtx(ctx, kind, fullPath)
			gullCtx := http.NewHttpContext(ctx.Writer, ctx.Request, &proxy{ctx: ctx})
			_, _ = e(gullCtx, nil)
			gullCtx.Push()
		})
	case endpoint.MethodPost:
		r.ginEngine.POST(fullPath.Path(), func(ctx *gin.Context) {
			initCtx(ctx, kind, fullPath)
			gullCtx := http.NewHttpContext(ctx.Writer, ctx.Request, &proxy{ctx: ctx})
			_, _ = e(gullCtx, nil)
			gullCtx.Push()
		})
	case endpoint.MethodPut:
		r.ginEngine.PUT(fullPath.Path(), func(ctx *gin.Context) {
			initCtx(ctx, kind, fullPath)
			gullCtx := http.NewHttpContext(ctx.Writer, ctx.Request, &proxy{ctx: ctx})
			_, _ = e(gullCtx, nil)
			gullCtx.Push()
		})
	case endpoint.MethodPatch:
		r.ginEngine.PATCH(fullPath.Path(), func(ctx *gin.Context) {
			initCtx(ctx, kind, fullPath)
			gullCtx := http.NewHttpContext(ctx.Writer, ctx.Request, &proxy{ctx: ctx})
			_, _ = e(gullCtx, nil)
			gullCtx.Push()
		})
	case endpoint.MethodDelete:
		r.ginEngine.DELETE(fullPath.Path(), func(ctx *gin.Context) {
			initCtx(ctx, kind, fullPath)
			gullCtx := http.NewHttpContext(ctx.Writer, ctx.Request, &proxy{ctx: ctx})
			_, _ = e(gullCtx, nil)
			gullCtx.Push()
		})
	case endpoint.MethodOptions:
		r.ginEngine.OPTIONS(fullPath.Path(), func(ctx *gin.Context) {
			initCtx(ctx, kind, fullPath)
			gullCtx := http.NewHttpContext(ctx.Writer, ctx.Request, &proxy{ctx: ctx})
			_, _ = e(gullCtx, nil)
			gullCtx.Push()
		})
	default:

	}
}

func initCtx(ctx *gin.Context, kind endpoint.Kind, fullPath endpoint.Path) {
	endpointCtxVal := endpoint.NewTransporter(kind, endpoint.NewHeader(ctx.Request.Header), fullPath.FullPath(), fullPath.Method())
	ctx.Request = ctx.Request.WithContext(endpoint.WithContext(ctx.Request.Context(), endpointCtxVal))
}
