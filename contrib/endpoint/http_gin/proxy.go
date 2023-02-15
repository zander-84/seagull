package http_gin

import "github.com/gin-gonic/gin"

type proxy struct {
	ctx *gin.Context
}

func (g *proxy) Param(key string) string {
	return g.ctx.Param(key)
}
