package ginx

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "reflect"
)

// Handler a group of apis, support auto register routes to gin.RouterGroup
type Handler interface {
    // RegisterRoute register route internally
    RegisterRoute(g *gin.RouterGroup)
}

// Request a request from remote client
type Request interface{}

// ValidateableRequest a request should be validated by call itself Validate function mannually
type ValidateableRequest interface {
    Validate() error
}

// NewRequest create a new request by the given request
func NewRequest(r Request) interface{} {
    if r == nil {
        return nil
    }
    rt := reflect.TypeOf(r)
    if rt.Kind() == reflect.Ptr {
        rt = rt.Elem()
    }
    return reflect.New(rt).Interface()
}

// Response any response send to client
type Response interface{}

// ApiResponser a responser was used to show request result to client
type ApiResponser interface {
    // Response a common response
    Response(c *gin.Context, code int, v interface{})
    // Success show a success response
    Success(c *gin.Context, v interface{})
    // Fail show a fail response
    Fail(c *gin.Context, v interface{})
}

// DefaultApiResponser a default responser implements responser interface
type DefaultApiResponser struct{}

func (r DefaultApiResponser) Response(c *gin.Context, code int, v interface{}) {
    c.JSON(code, v)
    c.Abort()
}

func (r DefaultApiResponser) Success(c *gin.Context, v interface{}) {
    c.JSON(http.StatusOK, v)
    c.Abort()
}

func (r DefaultApiResponser) Fail(c *gin.Context, v interface{}) {
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

// Logger define a log interface
type Logger interface {
    Debugf(format string, args ...interface{})
    Infof(format string, args ...interface{})
    Errorf(format string, args ...interface{})
}

type defaultLogger struct {
}

func (l *defaultLogger) Debugf(format string, args ...interface{}) {
    return
}

func (l *defaultLogger) Infof(format string, args ...interface{}) {
    return
}

func (l *defaultLogger) Errorf(format string, args ...interface{}) {
    return
}
