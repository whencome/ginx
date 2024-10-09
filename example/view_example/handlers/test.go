package handlers

import (
    "github.com/gin-gonic/gin"
    "github.com/whencome/ginx"
    "log"
    "view_example/reqs"
)

// Test handler
type Test struct{}

// Test handler func
func (t Test) Test(c *gin.Context, p *ginx.Page, r ginx.Request) error {
    req := r.(*reqs.TestRequest)
    log.Printf("request: %+v\n", req)
    if req.Name == "QUIT" {
        panic("test panic from view demo")
    }
    p.SetTitle("欢迎 " + req.Name)
    p.AddData("Name", req.Name)
    return nil
}
