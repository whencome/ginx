package main

import (
	"github.com/gin-gonic/gin"
	"github.com/whencome/ginx"
	"github.com/whencome/ginx/server"
)

var svr *server.HTTPServer

func main() {
	opts := &server.Options{
		Port: 8911,
	}
	svr = server.New(opts)
	svr.Init(initRoutes)
	// NOTE: server run in non-block mode
	svr.Start()
	// call Wait to prevent main goroutine exit immediately
	svr.Wait()
}

func initRoutes(r *gin.Engine) {
	r.GET("/greet", ginx.NewApi(GreetRequest{}, GreetLogic))
}
