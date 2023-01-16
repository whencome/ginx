package main

import (
    "github.com/gin-gonic/gin"
    "github.com/whencome/ginx"
    "github.com/whencome/ginx/server"
    "github.com/whencome/ginx/view"
    "log"
    "os"
    "path/filepath"
    "view_example/handlers"
    "view_example/reqs"
)

var svr *server.HTTPServer

func main() {
    // init view
    initView()
    // create && init http server
    opts := &server.Options{
        Port: 8912,
        Mode: server.ModeDebug,
    }
    svr = server.New(opts)
    svr.Init(initRoutes)
    // NOTE: server run in block mode
    if err := svr.Run(); err != nil {
        log.Printf("run server failed: %s\n", err)
        return
    }
}

// 初始化View
func initView() {
    // 获取工作目录
    wd, err := os.Getwd()
    if err != nil {
        wd = "./"
    }
    // 组装工作目录
    viewDir := filepath.Join(wd, "view")
    log.Printf("view dir: %s", viewDir)
    // 设置view
    view.SetTplDir(viewDir)
    view.AddTplFiles("template/public", "template/navbar", "template/error")
}

func initRoutes(r *gin.Engine) {
    // init routes
    r.GET("/test", ginx.NewPageHandler(reqs.TestRequest{}, "test/test", new(handlers.Test).Test))
}
