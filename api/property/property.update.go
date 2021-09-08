package property

import (
	"etender/api/handler"
	"etender/mysql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateDivision(c *gin.Context) {
}

func UpdateSSG(c *gin.Context) {
	ssgId := c.Param("ssgId")
	if ssgId != "" {
		var updateData SSG
		err := c.BindJSON(&updateData)
		fmt.Println(updateData)
		if err != nil {
			handler.ErrorHandler(c, http.StatusBadRequest, "Send Data in proper format", fmt.Errorf("invalid data"))
		}
		mySql := mysql.MysqlDB()
		defer mySql.Close()
		sqlStatement := "UPDATE ssg SET reserveprice = " + updateData.ReservePrice + " , emd = " + updateData.EMD + " WHERE ssgId =?;"

		query, err := mySql.Prepare(sqlStatement)
		if err != nil {
			handler.ErrorHandler(c, http.StatusBadRequest, "Could not update data", err)
		}
		res, err := query.Exec(ssgId)
		if err != nil {
			handler.ErrorHandler(c, http.StatusBadRequest, "Could not update data", err)
		} else {
			handler.SuccessHandler(c, http.StatusOK, "Succesfully Updated", res)
		}
	} else {
		handler.ErrorHandler(c, http.StatusBadRequest, "Invalid Param", fmt.Errorf(""))
	}

}

func UpdateFRE(c *gin.Context) {

}
