package types

import "github.com/gin-gonic/gin"

// Handler a group of apis, support auto register routes to gin.RouterGroup
type Handler interface {
    // RegisterRoute register route internally
    RegisterRoute(g *gin.RouterGroup)
}
