package main

import (
    "errors"
    "fmt"
    "log"

    "github.com/gin-gonic/gin"
    "github.com/whencome/ginx"
    "github.com/whencome/xlog"
)

var svr *ginx.HTTPServer

func main() {
    ginx.UseLogger(xlog.Use("default"))
    // register api
    ginx.UseApiResponser(new(ApiResponser))
    // run server
    opts := &ginx.ServerOptions{
        Port: 8915,
        Mode: ginx.ModeDebug,
    }
    svr = ginx.NewServer(opts)
    svr.PostInit(func(r *gin.Engine) error {
        initRoutes(r)
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
    r.Use(ginx.NewHandler(testMiddleware))
    r.GET("/greet", ginx.NewApiHandler(GreetRequest{}, GreetLogic))
}

func testMiddleware(c *gin.Context) error {
    fmt.Println("--------------$$ print from test middleware $$-----------------")
    if c.DefaultQuery("name", "whencome") == "whencome" {
        return errors.New("whencome not allowed here")
    }
    return nil
}
