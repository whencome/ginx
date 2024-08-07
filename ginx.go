package ginx

import (
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/whencome/ginx/log"
    "github.com/whencome/ginx/validator"
)

const (
    // the key for request cache
    requestKey = "__ginx_request__"
)

// UseLogger register a global logger
func UseLogger(l log.Logger) {
    log.Use(l)
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

// RequestParams get current request content
// 此方法用于在日志记录时获取请求内容
func RequestParams(c *gin.Context) interface{} {
    if v, ok := c.Get(requestKey); ok {
        return v
    }
    if c.Request.Method == http.MethodGet {
        return c.Request.Form
    }
    contentType := c.ContentType()
    if contentType == "application/x-www-form-urlencoded" {
        return c.Request.PostForm
    }
    if contentType == "multipart/form-data" {
        return c.Request.MultipartForm.Value
    }
    return nil
}

// HandlerFunc the logic to handle the api request
type HandlerFunc func(c *gin.Context) error

// NewHandler create a new gin.HandlerFunc, with no request, it's often been used to wrap a middleware
func NewHandler(f HandlerFunc) gin.HandlerFunc {
    return func(c *gin.Context) {
        if f == nil {
            return
        }
        err := f(c)
        if err != nil {
            getApiResponser().Fail(c, err)
            c.Abort()
            return
        }
    }
}

// ApiHandlerFunc the logic to handle the api request
type ApiHandlerFunc func(c *gin.Context, r Request) (Response, error)

// ApiMiddleware add middleware support to each single handler
type ApiMiddleware func(ApiHandlerFunc) ApiHandlerFunc

func apiMiddlewareChain(ms ...ApiMiddleware) ApiMiddleware {
    return func(next ApiHandlerFunc) ApiHandlerFunc {
        for i := len(ms) - 1; i >= 0; i-- {
            next = ms[i](next)
        }
        return next
    }
}

// NewApiHandler create a new gin.HandlerFunc
func NewApiHandler(r Request, f ApiHandlerFunc, ms ...ApiMiddleware) gin.HandlerFunc {
    return func(c *gin.Context) {
        if f == nil {
            getApiResponser().Response(c, http.StatusNotImplemented, "service not implemented")
            c.Abort()
            return
        }
        // parse && validate request
        var req Request
        if r != nil {
            req = NewRequest(r)
            if err := c.ShouldBind(req); err != nil {
                getApiResponser().Response(c, http.StatusBadRequest, validator.Error(err))
                c.Abort()
                return
            }
            c.Set(requestKey, req)
            if vr, ok := req.(ValidatableRequest); ok {
                if err := vr.Validate(); err != nil {
                    getApiResponser().Response(c, http.StatusBadRequest, err)
                    c.Abort()
                    return
                }
            }
        }
        // execute chain call
        var resp Response
        var err error
        if len(ms) > 0 {
            resp, err = apiMiddlewareChain(ms...)(f)(c, req)
        } else {
            resp, err = f(c, req)
        }
        if err != nil {
            getApiResponser().Fail(c, err)
            c.Abort()
            return
        }
        // if the response is nil, then won't use the responser to make a success response.
        if resp != nil {
            getApiResponser().Success(c, resp)
        }
    }
}

// PageHandlerFunc the logic to handle the page request
type PageHandlerFunc func(c *gin.Context, p *Page, r Request) error

// PageMiddleware add middleware support to each single handler
type PageMiddleware func(PageHandlerFunc) PageHandlerFunc

func pageMiddlewareChain(ms ...PageMiddleware) PageMiddleware {
    return func(next PageHandlerFunc) PageHandlerFunc {
        for i := len(ms) - 1; i >= 0; i-- {
            next = ms[i](next)
        }
        return next
    }
}

// NewPageHandler 创建一个页面处理方法
// t - template of current page
// r - request
// f - handler func
// ms - middleware list, allow empty
func NewPageHandler(v *View, t string, r Request, f PageHandlerFunc, ms ...PageMiddleware) gin.HandlerFunc {
    return func(c *gin.Context) {
        p := NewPage(c, v, t)
        if f == nil {
            _ = p.ShowWithError("service not implemented")
            c.Abort()
            return
        }
        // parse && validate request
        var req Request
        if r != nil {
            req = NewRequest(r)
            if err := c.ShouldBind(req); err != nil {
                _ = p.ShowWithError(validator.Error(err))
                c.Abort()
                return
            }
            c.Set(requestKey, req)
            if vr, ok := req.(ValidatableRequest); ok {
                if err := vr.Validate(); err != nil {
                    _ = p.ShowWithError(err)
                    c.Abort()
                    return
                }
            }
        }
        var err error
        if len(ms) > 0 {
            err = pageMiddlewareChain(ms...)(f)(c, p, req)
        } else {
            err = f(c, p, req)
        }
        if err != nil {
            _ = p.ShowWithError(err)
            c.Abort()
            return
        }
        _ = p.Show()
    }
}

// Wait 信号监听，当监听到退出信号时，执行资源释放函数，并退出程序
// f 程序退出前的资源释放方法
func Wait(releaseFunc func()) {
    sigChan := make(chan os.Signal)
    signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
    for {
        s := <-sigChan
        switch s {
        case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP:
            log.Infof("recv exit signal...")
            // 释放相关资源
            if releaseFunc != nil {
                log.Infof("execute release func...")
                releaseFunc()
            }
            // 等待1秒再退出
            time.Sleep(1 * time.Second)
            log.Infof("exit app...")
            os.Exit(0)
            return
        }
    }
}
