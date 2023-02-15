package http_router

import (
	"github.com/julienschmidt/httprouter"
)

type proxy struct {
	params httprouter.Params
}

func (p *proxy) Param(key string) string {
	return p.params.ByName(key)
}
