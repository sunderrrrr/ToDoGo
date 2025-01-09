package handler

import (
	"github.com/gin-gonic/gin"
	"log"
)

type error struct {
	msg string `json:"msg"`
}

func newErrorResponse(c *gin.Context, statusCode int, msg string) {
	log.Println(msg)
	c.AbortWithStatusJSON(statusCode, error{msg: msg})
}
