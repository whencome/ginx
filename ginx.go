package ginx

import (
    "github.com/gin-gonic/gin"
    "github.com/whencome/ginx/validator"
    "net/http"
)

// global logger
var logger Logger = new(defaultLogger)

// UseLogger register a global logger
func UseLogger(l Logger) {
    if l != nil {
        logger = l
    }
}

// apiResponser responser for api
var apiResponser ApiResponser

// UseApiResponser register a customized responser
func UseApiResponser(r ApiResponser) {
    apiResponser = r
}

// getApiResponser get the current available responser, if no customized resposer registered, a default api responser will be returned
func getApiResponser() ApiResponser {
    if apiResponser == nil {
        return new(DefaultApiResponser)
    }
    return apiResponser
}

// ApiHandlerFunc the logic to handle the api request
type ApiHandlerFunc func(c *gin.Context, r Request) (Response, error)

// NewApiHandler create a new gin.HandlerFunc
func NewApiHandler(r Request, f ApiHandlerFunc) gin.HandlerFunc {
    return func(c *gin.Context) {
        if f == nil {
            getApiResponser().Response(c, http.StatusNotImplemented, "service not implemented")
            return
        }
        // parse && validate request
        var req Request
        if r != nil {
            req = NewRequest(r)
            if err := c.ShouldBind(req); err != nil {
                getApiResponser().Response(c, http.StatusBadRequest, validator.Error(err))
                return
            }
            if vr, ok := req.(ValidateableRequest); ok {
                if err := vr.Validate(); err != nil {
                    getApiResponser().Response(c, http.StatusBadRequest, err)
                    return
                }
            }
        }
        resp, err := f(c, req)
        if err != nil {
            getApiResponser().Fail(c, err)
            return
        }
        getApiResponser().Success(c, resp)
        return
    }
}

// PageHandlerFunc the logic to handle the page request
type PageHandlerFunc func(c *gin.Context, p *Page, r Request) error

// NewPageHandler 创建一个页面处理方法
// t - template of current page
// r - request
// f - handler func
func NewPageHandler(v *View, t string, r Request, f PageHandlerFunc) gin.HandlerFunc {
    return func(c *gin.Context) {
        p := NewPage(c, v, t)
        if f == nil {
            _ = p.ShowWithError("service not implemented")
            return
        }
        // parse && validate request
        var req Request
        if r != nil {
            req = NewRequest(r)
            if err := c.ShouldBind(req); err != nil {
                _ = p.ShowWithError(validator.Error(err))
                return
            }
            if vr, ok := req.(ValidateableRequest); ok {
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
