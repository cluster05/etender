package property

import "github.com/gin-gonic/gin"

func Routes(router *gin.Engine) {
	router.POST("/excel", SaveExcel)

	router.GET("/data/division", GetDivision)
	router.GET("/data/ssg/:divisionId", GetSSG)
	router.GET("/data/fre/:ssgid", GetFRE)
}
