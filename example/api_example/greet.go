package main

import (
    "errors"
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/whencome/ginx"
    "log"
)

type GreetRequest struct {
    Name string `form:"name" label:"姓名" binding:"required" binding:"required"`
}

func GreetLogic(c *gin.Context, r ginx.Request) (ginx.Response, error) {
    // a type convert was needed
    req := r.(*GreetRequest)
    return fmt.Sprintf("hello %s", req.Name), nil
}

func LogMiddleware(f ginx.ApiHandlerFunc) ginx.ApiHandlerFunc {
    return func(c *gin.Context, r ginx.Request) (ginx.Response, error) {
        log.Printf("[LogLogic] request: %+v\n", r)
        ret, err := f(c, r)
        log.Printf("[LogLogic] response: %+v; err: %v\n", ret, err)
        return ret, err
    }
}

func FilterMiddleware(f ginx.ApiHandlerFunc) ginx.ApiHandlerFunc {
    return func(c *gin.Context, r ginx.Request) (ginx.Response, error) {
        log.Printf("[FilterLogic] request: %+v\n", r)
        // a type convert was needed
        req, ok := r.(*GreetRequest)
        if !ok {
            c.Abort()
            return nil, errors.New("invalid request")
        }
        // check name
        if req.Name == "eric" {
            c.Abort()
            return nil, errors.New("you are in black list, access denied")
        }
        req.Name += "_filtered"
        return f(c, req)
    }
}
