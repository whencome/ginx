package types

import (
    "github.com/gin-gonic/gin"
    ut "github.com/go-playground/universal-translator"
    "github.com/go-playground/validator/v10"
)

// Handler a group of apis, support auto register routes to gin.RouterGroup
type Handler interface {
    // RegisterRoute register route internally
    RegisterRoute(g *gin.RouterGroup)
}

// ErrorTranslator 错误解释器，用于多语言环境的错误处理
type ErrorTranslator interface {
    RegisterTranslations(v *validator.Validate) (ut.Translator, error)
}
