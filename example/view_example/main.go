package main

import (
    "github.com/gin-gonic/gin"
    "github.com/whencome/ginx"
    "log"
    "os"
    "path/filepath"
    "view_example/handlers"
    "view_example/reqs"
)

var svr *ginx.HTTPServer

func main() {
    // create && init http server
    opts := &ginx.ServerOptions{
        Port: 8912,
        Mode: ginx.ModeDebug,
    }
    svr = ginx.NewServer(opts)
    svr.PreInit(initRoutes)
    // NOTE: server run in block mode
    if err := svr.Run(); err != nil {
        log.Printf("run server failed: %s\n", err)
        return
    }
}

// 初始化View
func initView() *ginx.View {
    // 获取工作目录
    wd, err := os.Getwd()
    if err != nil {
        wd = "./"
    }
    // 组装工作目录
    viewDir := filepath.Join(wd, "view")
    log.Printf("view dir: %s", viewDir)
    view := ginx.NewView(
        ginx.WithTplDir(viewDir),
        ginx.WithTplFiles("template/public", "template/navbar", "template/error"))
    return view
}

func initRoutes(r *gin.Engine) error {
    // init routes
    v := initView()
    r.GET("/test", ginx.NewPageHandler(v, "test/test", reqs.TestRequest{}, new(handlers.Test).Test))
    r.GET("/test_middleware", ginx.NewPageHandler(v, "test/test", reqs.TestRequest{}, new(handlers.Test).Test, LogMiddleware))
    return nil
}

func LogMiddleware(f ginx.PageHandlerFunc) ginx.PageHandlerFunc {
    return func(c *gin.Context, p *ginx.Page, r ginx.Request) error {
        log.Printf("[LogLogic] request: %+v\n", r)
        err := f(c, p, r)
        log.Printf("[LogLogic] err: %v\n", err)
        return err
    }
}
