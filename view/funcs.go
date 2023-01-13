package view

import (
    "github.com/gin-gonic/gin"
    "github.com/whencome/ginx/types"
)

// NewHandler 创建一个页面处理方法
func NewHandler(t string, r types.Request, f types.PageLogicFunc) gin.HandlerFunc {
    return func(c *gin.Context) {
        p := NewPage(c, t)
        if f == nil {
            p.ShowWithError("service not implemented")
            return
        }
        // parse && validate request
        var req types.Request
        if r != nil {
            req = types.NewRequest(r)
            if err := c.ShouldBind(req); err != nil {
                p.ShowWithError(err)
                return
            }
            if vr, ok := req.(types.ValidateableRequest); ok {
                if err := vr.Validate(); err != nil {
                    p.ShowWithError(err)
                    return
                }
            }
        }
        err := f(c, p, req)
        if err != nil {
            p.ShowWithError(err)
            return
        }
        p.Show()
        return
    }
}
