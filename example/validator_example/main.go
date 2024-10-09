package main

import (
    "github.com/gin-gonic/gin"
    "github.com/go-playground/locales/en"
    ut "github.com/go-playground/universal-translator"
    "github.com/go-playground/validator/v10"
    et "github.com/go-playground/validator/v10/translations/en"
    "github.com/whencome/ginx"
    v "github.com/whencome/ginx/validator"
    "log"
)

var svr *ginx.HTTPServer

// ErrTrans 实现一个自定义的解释器
type ErrTrans struct{}

func (t *ErrTrans) RegisterTranslations(v *validator.Validate) (ut.Translator, error) {
    translator := en.New()
    trans, _ := ut.New(translator, translator).GetTranslator("en")
    et.RegisterDefaultTranslations(v, trans)
    return trans, nil
}

func main() {
    // 设置错误分割符号
    v.SetErrSeparator("||")
    // 显示全部错误
    v.ShowFullError(true)
    // 使用自定义解释器
    v.UseTranslator(new(ErrTrans))
    // run server
    opts := &ginx.ServerOptions{
        Port: 8914,
        Mode: ginx.ModeDebug,
    }
    svr = ginx.NewServer(opts)
    svr.PostInit(func(r *gin.Engine) error {
        initRoutes(r)
        return nil
    })
    // NOTE: server run in non-block mode
    _, err := svr.Start()
    if err != nil {
        log.Printf("start server failed: %s\n", err)
        return
    }
    log.Printf("start server on port %d in %s mode", opts.Port, opts.Mode)
    // call Wait to prevent main goroutine exit immediately
    svr.Wait()
}

func initRoutes(r *gin.Engine) {
    r.GET("/greet", ginx.NewApiHandler(GreetRequest{}, GreetLogic))
}
