package property

import (
	"etender/api/handler"
	"etender/mysql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UpdateDivision(c *gin.Context) {
}

func UpdateSSG(c *gin.Context) {
	ssgId := c.Param("ssgId")
	if ssgId != "" {
		var updateData = struct {
			ReservePrice int
			EMD          int
		}{}

		err := c.BindJSON(&updateData)

		reservePrice := strconv.Itoa(updateData.ReservePrice)
		emd := strconv.Itoa(updateData.EMD)

		if err != nil {
			handler.ErrorHandler(c, http.StatusBadRequest, "Send Data in proper format", fmt.Errorf("invalid data"))
		}
		mySql := mysql.MysqlDB()
		defer mySql.Close()
		sqlStatement := "UPDATE ssg SET reserveprice = " + reservePrice + " , emd = " + emd + " WHERE ssgId = ?;"

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
