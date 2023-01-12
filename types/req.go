package types

import (
    "ginx/view"
    "github.com/gin-gonic/gin"
    "reflect"
)

// Request a request from remote client
type Request interface{}

// ValidateableRequest a request should be validated by call itself Validate function mannually
type ValidateableRequest interface {
    Validate() error
}

// Response any response send to client
type Response interface{}

// ApiLogicFunc the logic to handle the api request
type ApiLogicFunc func(c *gin.Context, r Request) (Response, error)

// PageLogicFunc the logic to handle the page request
type PageLogicFunc func(c *gin.Context, p *view.Page, r Request) error

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
