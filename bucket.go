package ginx

import (
	"github.com/gin-gonic/gin"
)

// bucketer a bucket interface
type bucketer interface {
	Register()
	Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup
}

// Bucket a bucket is a group of apis, like *gin.RouterGroup
type Bucket struct {
	routerGroup *gin.RouterGroup
	middlewares []gin.HandlerFunc
	handlers    []Handler
	buckets     []bucketer
}

// NewBucket create a new bucket
func NewBucket(r *gin.RouterGroup, handlers ...Handler) *Bucket {
	b := &Bucket{
		routerGroup: r,
		middlewares: make([]gin.HandlerFunc, 0),
		handlers:    make([]Handler, 0),
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
			h.RegisterRoute(b.routerGroup)
		}
	}
	if len(b.buckets) > 0 {
		for _, bb := range b.buckets {
			bb.Register()
		}
	}
}

func (b *Bucket) Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return b.routerGroup.Group(relativePath, handlers...)
}

func (b *Bucket) UseMiddlewares(ms ...gin.HandlerFunc) {
	if len(ms) == 0 {
		return
	}
	b.routerGroup.Use(ms...)
}

func (b *Bucket) AddHandler(h Handler) {
	if b.handlers == nil {
		b.handlers = make([]Handler, 0)
	}
	b.handlers = append(b.handlers, h)
}

func (b *Bucket) AddHandlers(hs []Handler) {
	if b.handlers == nil {
		b.handlers = make([]Handler, 0)
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
