package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/whencome/ginx"
)

var svr *ginx.HTTPServer

func main() {
	// register global middleware
	ginx.UseApiMiddleware(Recover)
	// register api
	ginx.UseApiResponser(new(ApiResponser))
	// run server
	opts := &ginx.ServerOptions{
		Port: 8911,
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

// LogAccess 一个测试中间件
func LogAccess(c *gin.Context) {
	fmt.Println(time.Now().String())
}

func initRoutes(r *gin.Engine) {
	r.GET("/greet", ginx.NewApiHandler(GreetRequest{}, GreetLogic))
	r.GET("/greet_middleware", ginx.NewApiHandler(GreetRequest{}, GreetLogic, LogMiddleware, FilterMiddleware))
	r.GET("/greet/sayhi", ginx.NewApiHandler(SayHiRequest{}, SayHiLogic))
	r.GET("/time", ginx.NewApiHandler(nil, TimeLogic))
}

// Recover a middleware to capture panics
func Recover(f ginx.ApiHandlerFunc) ginx.ApiHandlerFunc {
	return func(c *gin.Context, r ginx.Request) (ginx.Response, error) {
		var ret ginx.Response
		var err error
		defer func() {
			if e := recover(); e != nil {
				log.Printf("panic: %v", e)
				err = fmt.Errorf("panic: %s", e)
			}
		}()
		ret, err = f(c, r)
		return ret, err
	}
}
