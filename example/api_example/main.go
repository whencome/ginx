package main

import (
    "log"

    "github.com/gin-gonic/gin"
    "github.com/whencome/ginx"
    "github.com/whencome/ginx/server"
)

var svr *server.HTTPServer

func main() {
    // register api
    ginx.UseApiResponser(new(ApiResponser))
    // run server
    opts := &server.Options{
        Port: 8911,
        Mode: server.ModeDebug,
    }
    svr = server.New(opts)
    svr.PreInit(func(r *gin.Engine) error {
        initRoutes(r)
        log.Println("--------- pre init ---------")
        return nil
    })
    svr.PostInit(func(r *gin.Engine) error {
        log.Println("--------- post init ---------")
        return nil
    })
    svr.PreStop(func(r *gin.Engine) error {
        log.Println("--------- pre stop ---------")
        return nil
    })
    svr.PostStop(func(r *gin.Engine) error {
        log.Println("--------- post stop ---------")
        return nil
    })
    // NOTE: server run in non-block mode
    ok, err := svr.Start()
    if err != nil {
        log.Printf("start server failed: %s\n", err)
        return
    }
    log.Printf("start server => %v\n", ok)
    // call Wait to prevent main goroutine exit immediately
    svr.Wait()
}

func initRoutes(r *gin.Engine) {
    r.GET("/greet", ginx.NewApiHandler(GreetRequest{}, GreetLogic))
}
