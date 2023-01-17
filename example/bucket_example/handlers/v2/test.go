package v2

import (
    "bucket_example/reqs"
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/whencome/ginx/api"
    "github.com/whencome/ginx/types"
)

type Test struct{}

func (t *Test) RegisterRoute(g *gin.RouterGroup) {
    r := g.Group("test")
    r.GET("/hi", api.NewHandler(reqs.SayHiRequest{}, t.SayHi))
    r.GET("/hello", api.NewHandler(reqs.SayHelloRequest{}, t.SayHello))
}

func (t *Test) SayHi(c *gin.Context, r types.Request) (types.Response, error) {
    req := r.(*reqs.SayHiRequest)
    return fmt.Sprintf("[V2] say hi to %s", req.Name), nil
}

func (t *Test) SayHello(c *gin.Context, r types.Request) (types.Response, error) {
    req := r.(*reqs.SayHelloRequest)
    return fmt.Sprintf("[V2] say hello to %s", req.Name), nil
}
