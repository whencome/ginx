package main

import (
    "fmt"
    "github.com/whencome/ginx/types"

    "github.com/gin-gonic/gin"
)

type GreetRequest struct {
    Name string `form:"name" label:"姓名" binding:"required" binding:"required"`
}

func GreetLogic(c *gin.Context, r types.Request) (types.Response, error) {
    // a type convert was needed
    req := r.(*GreetRequest)
    return fmt.Sprintf("hello %s", req.Name), nil
}
