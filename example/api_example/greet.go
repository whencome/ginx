package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/whencome/ginx"
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

// @Summary 打招呼
// @Description 用于向指定的对象打招呼
// @Produce json
// @Markdown
// ### 测试内容
// * string: 打招呼结果
// * error: 错误信息
// @Markdown
// @Router	/greet [post]
func GreetLogic(c *gin.Context, r ginx.Request) (ginx.Response, error) {
	// a type convert was needed
	req := r.(*GreetRequest)
	if req.Name == "QUIT" || req.Name == "QUIT_filtered" {
		panic("test panic from greet")
	}
	return fmt.Sprintf("hello %s", req.Name), nil
}

// 向别人说好
// @Summary SayHi
// @Description SayHi测试，SayHi to everyone
// @Produce text
// @Tags 问候
// @Markdown
// ### 返回内容
//
// ```json
//
//	{
//	  "message": "hello %s"
//	}
//
// ```
// @Markdown
// @Router	/greet/sayhi [get]
func SayHiLogic(c *gin.Context, r ginx.Request) (ginx.Response, error) {
	// a type convert was needed
	req := r.(*GreetRequest)
	if req.Name == "QUIT" || req.Name == "QUIT_filtered" {
		panic("test panic from greet")
	}
	return fmt.Sprintf("hello %s", req.Name), nil
}

// 显示当前时间
// @Summary Show Time
// @Description 显示当前时间
// @Produce json
// @Param timezone string false 时区
// @Markdown
// ### 返回内容
//
// ```json
//
//	{
//	  "message": "2005-01-02"
//	}
//
// ```
//
// **返回值说明**
//
// | 字段 | 类型 | 说明 |
// | --- | --- | --- |
// | message | string | 当前时间 |
//
// @Markdown
// @Router	/time [get]
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
