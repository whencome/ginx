package view

import (
    "github.com/gin-gonic/gin"
    "github.com/whencome/ginx/types"
    "github.com/whencome/ginx/validator"
)

// HandlerFunc the logic to handle the page request
type HandlerFunc func(c *gin.Context, p *Page, r types.Request) error

// NewHandler 创建一个页面处理方法
// t - template of current page
// r - request
// f - handler func
func NewHandler(t string, r types.Request, f HandlerFunc) gin.HandlerFunc {
    return func(c *gin.Context) {
        p := NewPage(c, t)
        if f == nil {
            _ = p.ShowWithError("service not implemented")
            return
        }
        // parse && validate request
        var req types.Request
        if r != nil {
            req = types.NewRequest(r)
            if err := c.ShouldBind(req); err != nil {
                _ = p.ShowWithError(validator.Error(err))
                return
            }
            if vr, ok := req.(types.ValidateableRequest); ok {
                if err := vr.Validate(); err != nil {
                    _ = p.ShowWithError(err)
                    return
                }
            }
        }
        err := f(c, p, req)
        if err != nil {
            _ = p.ShowWithError(err)
            return
        }
        _ = p.Show()
        return
    }
}
