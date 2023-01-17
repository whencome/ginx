package bucket

import (
    "github.com/gin-gonic/gin"
    "github.com/whencome/ginx/types"
)

// bucketer a bucket interface
type bucketer interface {
    Register()
    Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup
}

// Bucket a bucket is a group of apis, like *gin.RouterGroup
type Bucket struct {
    route       *gin.RouterGroup
    middlewares []gin.HandlerFunc
    handlers    []types.Handler
    buckets     []bucketer
}

func New(r *gin.RouterGroup, handlers ...types.Handler) *Bucket {
    b := &Bucket{
        route:       r,
        middlewares: make([]gin.HandlerFunc, 0),
        handlers:    make([]types.Handler, 0),
        buckets:     make([]bucketer, 0),
    }
    if len(handlers) > 0 {
        b.handlers = append(b.handlers, handlers...)
    }
    return b
}

func (b *Bucket) Register() {
    if len(b.handlers) > 0 {
        for _, h := range b.handlers {
            h.RegisterRoute(b.route)
        }
    }
    if len(b.buckets) > 0 {
        for _, bb := range b.buckets {
            bb.Register()
        }
    }
}

func (b *Bucket) Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
    return b.route.Group(relativePath, handlers...)
}

func (b *Bucket) UseMiddlewares(ms ...gin.HandlerFunc) {
    if len(ms) == 0 {
        return
    }
    b.route.Use(ms...)
}

func (b *Bucket) AddHandler(h types.Handler) {
    if b.handlers == nil {
        b.handlers = make([]types.Handler, 0)
    }
    b.handlers = append(b.handlers, h)
}

func (b *Bucket) AddHandlers(hs []types.Handler) {
    if b.handlers == nil {
        b.handlers = make([]types.Handler, 0)
    }
    b.handlers = append(b.handlers, hs...)
}

func (b *Bucket) AddBucket(b1 bucketer) {
    if b.buckets == nil {
        b.buckets = make([]bucketer, 0)
    }
    b.buckets = append(b.buckets, b1)
}

func (b *Bucket) AddBuckets(bs []bucketer) {
    if b.buckets == nil {
        b.buckets = make([]bucketer, 0)
    }
    b.buckets = append(b.buckets, bs...)
}
