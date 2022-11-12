package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/whencome/ginx/api"
	"github.com/whencome/ginx/server"
	"github.com/whencome/ginx/types"
)

// RegiserApiResponser 注册API Responser
func RegiserApiResponser(r types.Responser) {
	api.RegisterResponser(r)
}

// NewApi create a new gin.HandlerFunc
func NewApi(r api.Request, l api.LogicFunc) gin.HandlerFunc {
	return api.Api(r, l)
}

// NewServer create a new http server
func NewServer(opts *server.Options) *server.HTTPServer {
	return server.New(opts)
}
