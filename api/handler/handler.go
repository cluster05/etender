package handler

import (
	"github.com/gin-gonic/gin"
)

func ErrorHandler(c *gin.Context, status int, message string, err error) {
	c.JSON(status, gin.H{
		"message": message,
		"err":     err.Error(),
	})
}

func SuccessHandler(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, gin.H{
		"message":  message,
		"response": data,
	})
}
