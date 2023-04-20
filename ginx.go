package ginx

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/whencome/ginx/validator"
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

// NewApiHandler create a new gin.HandlerFunc
func NewApiHandler(r Request, f ApiHandlerFunc) gin.HandlerFunc {
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
			if vr, ok := req.(ValidateableRequest); ok {
				if err := vr.Validate(); err != nil {
					getApiResponser().Response(c, http.StatusBadRequest, err)
					c.Abort()
					return
				}
			}
		}
		resp, err := f(c, req)
		if err != nil {
			getApiResponser().Fail(c, err)
			c.Abort()
			return
		}
		getApiResponser().Success(c, resp)
		return
	}
}

// RawApiHandlerFunc defines a function that won't handle the corrent response, you should handle it by yourself,
// this function will handle the error so that make sure we handle the error in a unified way
type RawApiHandlerFunc func(c *gin.Context, r Request) error

// NewRawApiHandler create a new gin.HandlerFunc by a RawApiHandlerFunc
func NewRawApiHandler(r Request, f RawApiHandlerFunc) gin.HandlerFunc {
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
			if vr, ok := req.(ValidateableRequest); ok {
				if err := vr.Validate(); err != nil {
					getApiResponser().Response(c, http.StatusBadRequest, err)
					c.Abort()
					return
				}
			}
		}
		err := f(c, req)
		if err != nil {
			getApiResponser().Fail(c, err)
			c.Abort()
			return
		}
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
			if vr, ok := req.(ValidateableRequest); ok {
				if err := vr.Validate(); err != nil {
					_ = p.ShowWithError(err)
					c.Abort()
					return
				}
			}
		}
		err := f(c, p, req)
		if err != nil {
			_ = p.ShowWithError(err)
			c.Abort()
			return
		}
		_ = p.Show()
		return
	}
}
