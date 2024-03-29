package main

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/whencome/ginx"
)

type GreetRequest struct {
    Name string `form:"name" label:"姓名" binding:"required" binding:"required"`
}

func GreetLogic(c *gin.Context, r ginx.Request) (ginx.Response, error) {
    // a type convert was needed
    req := r.(*GreetRequest)
    return fmt.Sprintf("hello %s", req.Name), nil
}
