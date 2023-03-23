package v3

import (
    "bucket_example/reqs"
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/whencome/ginx"
)

type Test struct{}

func (t *Test) RegisterRoute(g *gin.RouterGroup) {
    r := g.Group("test")
    r.GET("/hi", ginx.NewApiHandler(reqs.SayHiRequest{}, t.SayHi))
    r.GET("/hello", ginx.NewApiHandler(reqs.SayHelloRequest{}, t.SayHello))
}

func (t *Test) SayHi(c *gin.Context, r ginx.Request) (ginx.Response, error) {
    req := r.(*reqs.SayHiRequest)
    return fmt.Sprintf("[V3] say hi to %s", req.Name), nil
}

func (t *Test) SayHello(c *gin.Context, r ginx.Request) (ginx.Response, error) {
    req := r.(*reqs.SayHelloRequest)
    return fmt.Sprintf("[V3] say hello to %s", req.Name), nil
}
