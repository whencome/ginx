package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/whencome/ginx/api"
)

type GreetRequest struct {
	Name string `form:"name" label:"姓名" binding:"required"`
}

func GreetLogic(c *gin.Context, r api.Request) (api.Response, error) {
	// a type convert was needed
	req := r.(*GreetRequest)
	return fmt.Sprintf("hello %s", req.Name), nil
}
