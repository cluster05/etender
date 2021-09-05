package auth

import (
	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {
	router.POST("/auth/register", AuthRegister)
	router.POST("/auth/login", AuthLogin)
}
