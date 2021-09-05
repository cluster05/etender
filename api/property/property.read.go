package property

import (
	"etender/api/handler"
	"etender/mysql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetDivision(c *gin.Context) {
	mySql := mysql.MysqlDB()
	stmt, err := mySql.Query("SELECT divisionId,name FROM division")

	var Divisions []DivisionView

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Query Failed",
			"err":     err.Error(),
		})
	}
	defer stmt.Close()

	for stmt.Next() {
		var result DivisionView
		err := stmt.Scan(&result.DivisionID, &result.Name)
		if err != nil {
			fmt.Printf("[GetData] Error Scssanning Data %v\n", err.Error())
			handler.ErrorHandler(c, "Query Failed", http.StatusInternalServerError, err.Error())
		}
		Divisions = append(Divisions, result)
	}

	handler.SuccessHandler(c, http.StatusOK, Divisions)

	// write else block and here

	defer mySql.Close()
}

func GetSSG(c *gin.Context) {

	queryData := c.Param("divisionId")

	if queryData != "" {
		mySql := mysql.MysqlDB()
		stmt, err := mySql.Query("SELECT ssgid,station,sector,pgroup FROM ssg WHERE divisionId = ?", queryData)

		if err != nil {
			handler.ErrorHandler(c, "Query Failed", http.StatusInternalServerError, err.Error())
		}
		defer stmt.Close()

		var Ssgs []SSG

		for stmt.Next() {
			var result SSG
			if err := stmt.Scan(&result.SSGId, &result.Station, &result.Sector,
				&result.Pgroup); err != nil {
				handler.ErrorHandler(c, "Query Failed", http.StatusInternalServerError, err.Error())
			}
			Ssgs = append(Ssgs, result)
		}
		handler.SuccessHandler(c, http.StatusOK, Ssgs)

		defer mySql.Close()
	} else {
		handler.ErrorHandler(c, "Pass Query", http.StatusBadRequest, "Param Not Passed")
	}
	// write else block and here

}
func GetFRE(c *gin.Context) {
	mySql := mysql.MysqlDB()
	queryData := c.Param("ssgid")

	stmt, err := mySql.Query("SELECT freid,flatno,reserveprice,emd FROM fre WHERE ssgId = ?", queryData)

	if queryData != "" {
		if err != nil {
			log.Println(err)
			handler.ErrorHandler(c, "Query Failed", http.StatusInternalServerError, err.Error())
		}
		defer stmt.Close()

		var Fres []FRE

		for stmt.Next() {
			var result FRE
			err = stmt.Scan(&result.FREId, &result.FlatNo, &result.ReservePrice, &result.EMD)
			if err != nil {
				handler.ErrorHandler(c, "Query Failed", http.StatusInternalServerError, err.Error())
			}
			Fres = append(Fres, result)
		}
		handler.SuccessHandler(c, http.StatusOK, Fres)
	} else {
		handler.ErrorHandler(c, "Pass Query", http.StatusBadRequest, "Param Not Passed")
	}
	// write else block and here
	defer mySql.Close()
}
