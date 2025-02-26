package handler

import (
	"github.com/gin-gonic/gin"
	"log"
)

type errorResponse struct {
	msg string `json:"msg"`
}

func newErrorResponse(c *gin.Context, statusCode int, msg string) {
	log.Println("response.go: " + msg)
	c.AbortWithStatusJSON(statusCode, errorResponse{msg: msg})
}
