package property

import "github.com/gin-gonic/gin"

func Routes(router *gin.Engine) {
	router.POST("/excel", SaveExcel)
}
