package main

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/whencome/ginx"
)

type GreetRequest struct {
    Greet string `form:"greet" label:"greeting word" binding:"required" binding:"required"`
    Name  string `form:"name" label:"name" binding:"required" binding:"required"`
}

func GreetLogic(c *gin.Context, r ginx.Request) (ginx.Response, error) {
    // a type convert was needed
    req := r.(*GreetRequest)
    return fmt.Sprintf("%s %s", req.Greet, req.Name), nil
}
