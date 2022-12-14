package types

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

// Responser a responser was used to show request result to client
type Responser interface {
    // Response a common response
    Response(c *gin.Context, code int, v interface{})
    // Success show a success response
    Success(c *gin.Context, v interface{})
    // Fail show a fail response
    Fail(c *gin.Context, v interface{})
}

// DefaultResponser a default responser implements responser interface
type DefaultResponser struct{}

func (r DefaultResponser) Response(c *gin.Context, code int, v interface{}) {
    c.JSON(code, v)
    c.Abort()
}

func (r DefaultResponser) Success(c *gin.Context, v interface{}) {
    c.JSON(http.StatusOK, v)
    c.Abort()
}

func (r DefaultResponser) Fail(c *gin.Context, v interface{}) {
    if e, ok := v.(ApiError); ok {
        c.JSON(e.Code(), v)
        c.Abort()
    }
    c.JSON(http.StatusBadRequest, v)
    c.Abort()
}

// ApiError the error should return an extra code that indicate which kind of error it is
type ApiError interface {
    error
    Code() int
}
