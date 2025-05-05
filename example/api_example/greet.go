package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/whencome/ginx"

	_ "github.com/whencome/ginx/example/api_example/response"
)

// GreetRequest 打招呼请求
type GreetRequest struct {
	// 姓名
	Name string `form:"name" label:"姓名" binding:"required"` // 打招呼对象
}

// SayHiRequest 问好请求
type SayHiRequest struct {
	// 姓名
	Name string `form:"name" label:"姓名" binding:"required"`
	// 时间
	Time string `form:"time" label:"时间"`
}

// GreetResponse 打招呼返回结果
type GreetResponse struct {
	Message string `json:"message"` // 打招呼结果
}

func GreetLogic(c *gin.Context, r ginx.Request) (ginx.Response, error) {
	// a type convert was needed
	req := r.(*GreetRequest)
	if req.Name == "QUIT" || req.Name == "QUIT_filtered" {
		panic("test panic from greet")
	}
	return fmt.Sprintf("hello %s", req.Name), nil
}

func SayHiLogic(c *gin.Context, r ginx.Request) (ginx.Response, error) {
	// a type convert was needed
	req := r.(*GreetRequest)
	if req.Name == "QUIT" || req.Name == "QUIT_filtered" {
		panic("test panic from greet")
	}
	return fmt.Sprintf("hello %s", req.Name), nil
}

func TimeLogic(c *gin.Context, r ginx.Request) (ginx.Response, error) {
	return "2005-01-02", nil
}

func LogMiddleware(f ginx.ApiHandlerFunc) ginx.ApiHandlerFunc {
	return func(c *gin.Context, r ginx.Request) (ginx.Response, error) {
		log.Printf("[LogLogic] request: %+v\n", r)
		ret, err := f(c, r)
		log.Printf("[LogLogic] response: %+v; err: %v\n", ret, err)
		return ret, err
	}
}

func FilterMiddleware(f ginx.ApiHandlerFunc) ginx.ApiHandlerFunc {
	return func(c *gin.Context, r ginx.Request) (ginx.Response, error) {
		log.Printf("[FilterLogic] request: %+v\n", r)
		// a type convert was needed
		req, ok := r.(*GreetRequest)
		if !ok {
			c.Abort()
			return nil, errors.New("invalid request")
		}
		// check name
		if req.Name == "eric" {
			c.Abort()
			return nil, errors.New("you are in black list, access denied")
		}
		req.Name += "_filtered"
		return f(c, req)
	}
}
