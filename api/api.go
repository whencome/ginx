package api

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/whencome/ginx/types"
)

// Request a request from remote client
type Request interface{}

// ValidateableRequest a request should be validated by call itself Validate function mannually
type ValidateableRequest interface {
	Validate() error
}

// Response any response send to client
type Response interface{}

// LogicFunc the logic to handle the request
type LogicFunc func(c *gin.Context, r Request) (Response, error)

// responser response result to client
var responser types.Responser

func getResponser() types.Responser {
	if responser == nil {
		return new(types.DefaultResponser)
	}
	return responser
}

// newRequest create a new request by the given request
func newRequest(r Request) interface{} {
	if r == nil {
		return nil
	}
	rt := reflect.TypeOf(r)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	return reflect.New(rt).Interface()
}

// newHandler create a new gin.HandlerFunc
func newHandler(r Request, f LogicFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if f == nil {
			getResponser().Response(c, http.StatusNotImplemented, "service not implemented")
			return
		}
		// parse && validate request
		var req Request
		if r != nil {
			req = newRequest(r)
			if err := c.ShouldBind(req); err != nil {
				getResponser().Fail(c, err)
				return
			}
			if vr, ok := req.(ValidateableRequest); ok {
				if err := vr.Validate(); err != nil {
					getResponser().Fail(c, err)
					return
				}
			}
		}
		resp, err := f(c, req)
		if err != nil {
			getResponser().Fail(c, err)
			return
		}
		getResponser().Success(c, resp)
		return
	}
}

// RegisterResponser register a customized responser to show result
func RegisterResponser(r types.Responser) {
	responser = r
}

// Api create a new api
func Api(r Request, f LogicFunc) gin.HandlerFunc {
	return newHandler(r, f)
}
