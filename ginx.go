package ginx

import (
    "ginx/api"
    "ginx/server"
    "ginx/types"
    "ginx/view"
    "github.com/gin-gonic/gin"
)

// UseApiResponser 注册API Responser
func UseApiResponser(r types.Responser) {
    api.UseResponser(r)
}

// NewServer create a new http server
func NewServer(opts *server.Options) *server.HTTPServer {
    return server.New(opts)
}

// NewApiHandler create a new gin.HandlerFunc
func NewApiHandler(r types.Request, l types.ApiLogicFunc) gin.HandlerFunc {
    return api.NewHandler(r, l)
}

// NewPageHandler 创建一个Page处理对象
func NewPageHandler(r types.Request, t string, l types.PageLogicFunc) gin.HandlerFunc {
    return view.NewHandler(t, r, l)
}
