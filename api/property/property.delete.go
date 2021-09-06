package property

import (
	"etender/api/handler"
	"etender/mysql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteDivision(c *gin.Context) {
	queryData := c.Param("divisionId")
	deleteData(c, queryData, "division", "divisionId")
}

func DeleteFRE(c *gin.Context) {
	queryData := c.Param("freId")
	deleteData(c, queryData, "fre", "freId")
}

func DeleteSSG(c *gin.Context) {
	queryData := c.Param("ssgId")
	deleteData(c, queryData, "ssg", "ssgId")
}

func deleteData(c *gin.Context, queryData, from, where string) {
	if queryData != "" {
		mySql := mysql.MysqlDB()
		defer mySql.Close()
		q := "DELETE FROM " + from + " Where " + where + "=?"
		delForm, err := mySql.Prepare(q)
		if err != nil {
			handler.ErrorHandler(c, http.StatusBadRequest, "Query Failed", err)
		}
		delForm.Exec(queryData)
	} else {
		handler.ErrorHandler(c, http.StatusBadRequest, "Could Not Delete", fmt.Errorf(""))
	}
}
