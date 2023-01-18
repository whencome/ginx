package ginx

import (
    "github.com/gin-gonic/gin"
    "github.com/whencome/ginx/api"
    "github.com/whencome/ginx/bucket"
    "github.com/whencome/ginx/server"
    "github.com/whencome/ginx/types"
    "github.com/whencome/ginx/view"
)

// UseApiResponser 注册API Responser
func UseApiResponser(r types.ApiResponser) {
    api.UseResponser(r)
}

// NewServer create a new http server
func NewServer(opts *server.Options) *server.HTTPServer {
    return server.New(opts)
}

// NewApiHandler create a new gin.HandlerFunc
func NewApiHandler(r types.Request, l api.HandlerFunc) gin.HandlerFunc {
    return api.NewHandler(r, l)
}

// NewPageHandler 创建一个Page处理对象
func NewPageHandler(r types.Request, t string, l view.HandlerFunc) gin.HandlerFunc {
    return view.NewHandler(t, r, l)
}

// NewBucket create a bucket
func NewBucket(r *gin.RouterGroup, handlers ...types.Handler) *bucket.Bucket {
    return bucket.New(r, handlers...)
}
