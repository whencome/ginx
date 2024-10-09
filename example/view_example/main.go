package main

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/whencome/ginx"
    "log"
    "os"
    "path/filepath"
    "runtime/debug"
    "view_example/handlers"
    "view_example/reqs"
)

var svr *ginx.HTTPServer

func main() {
    ginx.UsePageMiddleware(Recover)
    // create && init http server
    opts := &ginx.ServerOptions{
        Port: 8912,
        Mode: ginx.ModeDebug,
    }
    svr = ginx.NewServer(opts)
    svr.PostInit(initRoutes)
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
    r.Use(CustomRecovery)
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

// Recover a middleware to capture panics
func Recover(f ginx.PageHandlerFunc) ginx.PageHandlerFunc {
    return func(c *gin.Context, p *ginx.Page, r ginx.Request) error {
        var err error
        defer func() {
            if e := recover(); e != nil {
                log.Printf("panic: %v", e)
                err = fmt.Errorf("panic: %s", e)
            }
        }()
        err = f(c, p, r)
        return err
    }
}

// CustomRecovery 自定义 Recovery 中间件
func CustomRecovery(c *gin.Context) {
    defer func() {
        if err := recover(); err != nil {
            // 获取调用堆栈信息
            stackTrace := debug.Stack()
            log.Printf("[CustomRecovery]: %v\nStack Trace:\n%s\n", err, stackTrace)
            c.AbortWithStatus(500)
        }
    }()
    c.Next()
}
