package main

import (
    "bucket_example/handlers/v1"
    v2 "bucket_example/handlers/v2"
    v3 "bucket_example/handlers/v3"
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/whencome/ginx"
    "log"
)

var svr *ginx.HTTPServer

func main() {
    opts := &ginx.ServerOptions{
        Port: 8913,
        Mode: ginx.ModeDebug,
    }
    svr = ginx.NewServer(opts)
    svr.PreInit(initRoutes)
    if err := svr.Run(); err != nil {
        log.Printf("run server failed: %s\n", err)
        return
    }
}

func initRoutes(r *gin.Engine) error {
    v1g := r.Group("v1")
    v1Bucket := ginx.NewBucket(
        v1g,
        new(v1.Test),
        new(v1.Welcome))
    v1Bucket.Register()

    v2g := r.Group("v2")
    v2Bucket := ginx.NewBucket(
        v2g,
        new(v2.Test),
        new(v2.Welcome))

    v3g := v2g.Group("v3", func(context *gin.Context) {
        fmt.Println("-------v3 middleware-------")
    })
    v3Bucket := ginx.NewBucket(
        v3g,
        new(v3.Test),
        new(v3.Welcome))
    v2Bucket.AddBucket(v3Bucket)

    v2Bucket.Register()
    return nil
}
