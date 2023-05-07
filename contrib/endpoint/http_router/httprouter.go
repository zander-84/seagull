package http_router

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"github.com/zander-84/seagull/endpoint"
	"github.com/zander-84/seagull/transport"
	"github.com/zander-84/seagull/transport/http"
	http2 "net/http"
)

var _ http2.Handler = (*Router)(nil)

type Router struct {
	engine *httprouter.Router
}

func NewRouter(engine *httprouter.Router) *Router {
	router := new(Router)
	router.engine = engine
	return router
}

func (r *Router) ServeHTTP(res http2.ResponseWriter, req *http2.Request) {
	r.engine.ServeHTTP(res, req)
}
func (r *Router) Endpoint(kind transport.Kind, fullPath endpoint.Path, e endpoint.HandlerFunc) {
	switch fullPath.Method() {

	case transport.MethodGet:
		r.engine.GET(fullPath.Path(), func(writer http2.ResponseWriter, request *http2.Request, params httprouter.Params) {
			_, _ = e(initCtx(writer, request, kind, fullPath, params), nil)
		})
	case transport.MethodHead:
		r.engine.HEAD(fullPath.Path(), func(writer http2.ResponseWriter, request *http2.Request, params httprouter.Params) {
			_, _ = e(initCtx(writer, request, kind, fullPath, params), nil)
		})
	case transport.MethodPost:
		r.engine.POST(fullPath.Path(), func(writer http2.ResponseWriter, request *http2.Request, params httprouter.Params) {
			_, _ = e(initCtx(writer, request, kind, fullPath, params), nil)
		})
	case transport.MethodPut:
		r.engine.PUT(fullPath.Path(), func(writer http2.ResponseWriter, request *http2.Request, params httprouter.Params) {
			_, _ = e(initCtx(writer, request, kind, fullPath, params), nil)
		})
	case transport.MethodPatch:
		r.engine.PATCH(fullPath.Path(), func(writer http2.ResponseWriter, request *http2.Request, params httprouter.Params) {
			_, _ = e(initCtx(writer, request, kind, fullPath, params), nil)
		})
	case transport.MethodDelete:
		r.engine.DELETE(fullPath.Path(), func(writer http2.ResponseWriter, request *http2.Request, params httprouter.Params) {
			_, _ = e(initCtx(writer, request, kind, fullPath, params), nil)
		})
	case transport.MethodOptions:
		r.engine.OPTIONS(fullPath.Path(), func(writer http2.ResponseWriter, request *http2.Request, params httprouter.Params) {
			_, _ = e(initCtx(writer, request, kind, fullPath, params), nil)
		})
	default:
	}
}

func initCtx(writer http2.ResponseWriter, request *http2.Request, kind transport.Kind, fullPath endpoint.Path, params httprouter.Params) context.Context {
	endpointCtxVal := transport.NewTransporter(kind, transport.NewHeader(request.Header), fullPath.FullPath(), fullPath.Method())
	req := request.WithContext(transport.WithContext(request.Context(), endpointCtxVal))
	return http.NewHttpContext(writer, req, &proxy{params: params}, endpointCtxVal)
}
