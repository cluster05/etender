package property

import "github.com/gin-gonic/gin"

func Routes(router *gin.Engine) {
	router.POST("/excel", SaveExcel)

	router.GET("/data/division", GetDivision)
	router.GET("/data/ssg/:divisionId", GetSSG)
	router.GET("/data/fre/:ssgId", GetFRE)

	router.PATCH("/data/division/:divisionId", UpdateDivision)
	router.PATCH("/data/ssg/:ssgId", UpdateSSG)
	router.PATCH("/data/fre/:freId", UpdateFRE)

	router.DELETE("/data/division/:divisionId", DeleteDivision)
	router.DELETE("/data/ssg/:ssgId", DeleteSSG)
	router.DELETE("/data/fre/:freId", DeleteFRE)

}
