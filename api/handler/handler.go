package handler

import (
	"github.com/gin-gonic/gin"
)

func ErrorHandler(c *gin.Context, message string, status int, err string) {
	c.JSON(status, gin.H{
		"message": message,
		"err":     err,
	})
}

func SuccessHandler(c *gin.Context, status int, data interface{}) {
	c.JSON(status, gin.H{
		"data": data,
	})
}
