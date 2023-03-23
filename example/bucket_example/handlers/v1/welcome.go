package v1

import (
    "bucket_example/reqs"
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/whencome/ginx"
)

type Welcome struct{}

func (w *Welcome) RegisterRoute(g *gin.RouterGroup) {
    r := g.Group("welcome")
    r.GET("/greet", ginx.NewApiHandler(reqs.WelcomeRequest{}, w.Greet))
}

func (w *Welcome) Greet(c *gin.Context, r ginx.Request) (ginx.Response, error) {
    req := r.(*reqs.WelcomeRequest)
    return fmt.Sprintf("[V1] %s", req.Greet), nil
}
