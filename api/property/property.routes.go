package property

import "github.com/gin-gonic/gin"

func Routes(router *gin.Engine) {
	router.POST("/excel", SaveExcel)
	router.GET("/data/division", GetDivision)
	router.GET("/data/ssg/:division", GetSsg)
	router.GET("/data/fre/:ssgid", GetFre)
}
