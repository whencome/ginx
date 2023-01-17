package v2

import (
    "bucket_example/reqs"
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/whencome/ginx/api"
    "github.com/whencome/ginx/types"
)

type Welcome struct{}

func (w *Welcome) RegisterRoute(g *gin.RouterGroup) {
    r := g.Group("welcome")
    r.GET("/greet", api.NewHandler(reqs.WelcomeRequest{}, w.Greet))
}

func (w *Welcome) Greet(c *gin.Context, r types.Request) (types.Response, error) {
    req := r.(*reqs.WelcomeRequest)
    return fmt.Sprintf("[V2] %s", req.Greet), nil
}
