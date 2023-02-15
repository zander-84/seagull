package http_router

import (
	"github.com/julienschmidt/httprouter"
	"github.com/zander-84/seagull/endpoint"
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
func (r *Router) Endpoint(kind endpoint.Kind, fullPath endpoint.Path, e endpoint.HandlerFunc) {
	switch fullPath.Method() {

	case endpoint.MethodGet:
		r.engine.GET(fullPath.Path(), func(writer http2.ResponseWriter, request *http2.Request, params httprouter.Params) {
			request = setRequest(request, kind, fullPath)
			gullCtx := http.NewHttpContext(writer, request, &proxy{params: params})
			_, _ = e(gullCtx, nil)
		})
	case endpoint.MethodHead:
		r.engine.HEAD(fullPath.Path(), func(writer http2.ResponseWriter, request *http2.Request, params httprouter.Params) {
			request = setRequest(request, kind, fullPath)
			gullCtx := http.NewHttpContext(writer, request, &proxy{params: params})
			_, _ = e(gullCtx, nil)
		})
	case endpoint.MethodPost:
		r.engine.POST(fullPath.Path(), func(writer http2.ResponseWriter, request *http2.Request, params httprouter.Params) {
			request = setRequest(request, kind, fullPath)
			gullCtx := http.NewHttpContext(writer, request, &proxy{params: params})
			_, _ = e(gullCtx, nil)
		})
	case endpoint.MethodPut:
		r.engine.PUT(fullPath.Path(), func(writer http2.ResponseWriter, request *http2.Request, params httprouter.Params) {
			request = setRequest(request, kind, fullPath)
			gullCtx := http.NewHttpContext(writer, request, &proxy{params: params})
			_, _ = e(gullCtx, nil)
		})
	case endpoint.MethodPatch:
		r.engine.PATCH(fullPath.Path(), func(writer http2.ResponseWriter, request *http2.Request, params httprouter.Params) {
			request = setRequest(request, kind, fullPath)
			gullCtx := http.NewHttpContext(writer, request, &proxy{params: params})
			_, _ = e(gullCtx, nil)
		})
	case endpoint.MethodDelete:
		r.engine.DELETE(fullPath.Path(), func(writer http2.ResponseWriter, request *http2.Request, params httprouter.Params) {
			request = setRequest(request, kind, fullPath)
			gullCtx := http.NewHttpContext(writer, request, &proxy{params: params})
			_, _ = e(gullCtx, nil)
		})
	case endpoint.MethodOptions:
		r.engine.OPTIONS(fullPath.Path(), func(writer http2.ResponseWriter, request *http2.Request, params httprouter.Params) {
			request = setRequest(request, kind, fullPath)
			gullCtx := http.NewHttpContext(writer, request, &proxy{params: params})
			_, _ = e(gullCtx, nil)
		})
	default:
	}
}

func setRequest(request *http2.Request, kind endpoint.Kind, fullPath endpoint.Path) *http2.Request {
	endpointCtxVal := endpoint.NewTransporter(kind, endpoint.NewHeader(request.Header), fullPath.FullPath(), fullPath.Method())

	return request.WithContext(endpoint.WithContext(request.Context(), endpointCtxVal))
}
