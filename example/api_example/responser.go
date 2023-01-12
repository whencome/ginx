package main

import (
    "fmt"
    "net/http"

    "github.com/gin-gonic/gin"
)

type ApiResponseMessage struct {
    Code    int
    Message string
    Data    interface{}
}

type ApiResponser struct{}

func (r *ApiResponser) buildMsg(code int, v interface{}) ApiResponseMessage {
    msg := ApiResponseMessage{
        Code:    code,
        Message: "",
        Data:    nil,
    }
    // SUCCESS
    if code == http.StatusOK {
        msg.Message = "success"
        msg.Data = v
        return msg
    }
    // FAIL
    e, ok := v.(error)
    if ok {
        msg.Message = e.Error()
    } else {
        msg.Message = fmt.Sprintf("%s", v)
    }
    return msg
}

func (r *ApiResponser) Response(c *gin.Context, code int, v interface{}) {
    msg := r.buildMsg(code, v)
    c.JSON(code, msg)
    c.Abort()
}

func (r *ApiResponser) Success(c *gin.Context, v interface{}) {
    r.Response(c, http.StatusOK, v)
}

func (r *ApiResponser) Fail(c *gin.Context, v interface{}) {
    r.Response(c, http.StatusBadRequest, v)
}
