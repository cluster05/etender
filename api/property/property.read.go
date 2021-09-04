package property

import (
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
			fmt.Printf("[GetData] Error Scanning Data %v\n", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Query Failed",
				"err":     err.Error(),
			})
		}
		Divisions = append(Divisions, result)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": Divisions,
	})

	defer mySql.Close()
}

func GetSSG(c *gin.Context) {

	queryData := c.Param("divisionId")

	if queryData != "" {
		mySql := mysql.MysqlDB()
		stmt, err := mySql.Query("SELECT ssgid,station,sector,pgroup FROM ssg WHERE divisionId = ?", queryData)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Query Failed",
				"err":     err.Error(),
			})
		}
		defer stmt.Close()

		var Ssgs []SSG

		for stmt.Next() {
			var result SSG
			if err := stmt.Scan(&result.SSGId, &result.Station, &result.Sector,
				&result.Pgroup); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Query Failed",
					"err":     err.Error(),
				})
			}
			Ssgs = append(Ssgs, result)
		}

		c.JSON(http.StatusOK, gin.H{
			"data": Ssgs,
		})

		defer mySql.Close()
	}

}
func GetFRE(c *gin.Context) {
	mySql := mysql.MysqlDB()
	queryData := c.Param("ssgid")

	stmt, err := mySql.Query("SELECT freid,flatno,reserveprice,emd FROM fre WHERE ssgId = ?", queryData)

	if queryData != "" {
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Query Failed",
				"err":     err.Error(),
			})
		}
		defer stmt.Close()

		var Fres []FRE

		for stmt.Next() {
			var result FRE
			err = stmt.Scan(&result.FREId, &result.FlatNo, &result.ReservePrice, &result.EMD)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Query Failed",
					"err":     err.Error(),
				})
			}
			Fres = append(Fres, result)
		}

		c.JSON(http.StatusOK, gin.H{
			"data": Fres,
		})

	}
	defer mySql.Close()
}
