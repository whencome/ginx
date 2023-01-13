package api

import (
    "github.com/gin-gonic/gin"
    "github.com/whencome/ginx/types"
    "net/http"
)

// responser response result to client
var responser types.Responser

func getResponser() types.Responser {
    if responser == nil {
        return new(types.DefaultResponser)
    }
    return responser
}

// HandlerFunc the logic to handle the api request
type HandlerFunc func(c *gin.Context, r types.Request) (types.Response, error)

// NewHandler create a new gin.HandlerFunc
func NewHandler(r types.Request, f HandlerFunc) gin.HandlerFunc {
    return func(c *gin.Context) {
        if f == nil {
            getResponser().Response(c, http.StatusNotImplemented, "service not implemented")
            return
        }
        // parse && validate request
        var req types.Request
        if r != nil {
            req = types.NewRequest(r)
            if err := c.ShouldBind(req); err != nil {
                getResponser().Response(c, http.StatusBadRequest, err)
                return
            }
            if vr, ok := req.(types.ValidateableRequest); ok {
                if err := vr.Validate(); err != nil {
                    getResponser().Response(c, http.StatusBadRequest, err)
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

// UseResponser register a customized responser to show result
func UseResponser(r types.Responser) {
    responser = r
}
