package property

import (
	"etender/api/handler"
	"etender/mysql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetDivision(c *gin.Context) {
	mySql := mysql.MysqlDB()
	defer mySql.Close()
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
			handler.ErrorHandler(c, http.StatusInternalServerError, "Query Failed", err)
		}
		Divisions = append(Divisions, result)
	}

	handler.SuccessHandler(c, http.StatusOK, "success", Divisions)

	// write else block and here
}

func GetSSG(c *gin.Context) {
	queryData := c.Param("divisionId")

	if queryData != "" {
		mySql := mysql.MysqlDB()
		defer mySql.Close()
		stmt, err := mySql.Query("SELECT ssgid,station,sector,pgroup,reserveprice,emd FROM ssg WHERE divisionId = ?", queryData)

		if err != nil {
			handler.ErrorHandler(c, http.StatusInternalServerError, "Query Failed", err)
		}
		defer stmt.Close()
		var SSGMap = make(map[string]map[string][]SSG)
		var map1 = make(map[string][]SSG)
		for stmt.Next() {
			var result SSG
			if err := stmt.Scan(&result.SSGId, &result.Station, &result.Sector,
				&result.Pgroup, &result.ReservePrice, &result.EMD); err != nil {
				handler.ErrorHandler(c, http.StatusInternalServerError, "Query Failed", err)
			}
			var ssgIdPgroup SSG
			ssgIdPgroup.SSGId = result.SSGId
			ssgIdPgroup.Pgroup = result.Pgroup
			ssgIdPgroup.ReservePrice = result.ReservePrice
			ssgIdPgroup.EMD = result.EMD
			if _, ok := SSGMap[result.Station]; !ok {
				map1 = make(map[string][]SSG)
			}
			map1[result.Sector] = append(map1[result.Sector], ssgIdPgroup)
			SSGMap[result.Station] = map1
		}
		var tree []Tree
		for i, v := range SSGMap {
			var tree1 Tree
			tree1.Name = i
			tree1.EditNodeDisabled = true
			for i2, v2 := range v {
				var tree2 Tree
				tree2.Name = i2
				tree2.EditNodeDisabled = true

				for _, v3 := range v2 {
					var tree3 Tree

					reservePeice, _ := strconv.Atoi(v3.ReservePrice)
					emd, _ := strconv.Atoi(v3.EMD)
					tree3.Name = v3.Pgroup + "[RP] " + v3.ReservePrice + " [EMD] " + v3.EMD
					tree3.Text = v3.Pgroup
					tree3.SsgId = v3.SSGId
					tree3.ReservePrice = reservePeice
					tree3.EMD = emd
					tree3.Level = "ssg"
					tree3.EditNodeDisabled = false
					tree3.Children = []Tree{}
					tree2.Children = append(tree2.Children, tree3)
				}
				tree1.Children = append(tree1.Children, tree2)
			}
			tree = append(tree, tree1)
		}
		handler.SuccessHandler(c, http.StatusOK, "success", tree)
	} else {
		handler.ErrorHandler(c, http.StatusBadRequest, "Query Failed", fmt.Errorf(""))
	}
	// write else block and here

}
func GetFRE(c *gin.Context) {
	queryData := c.Param("ssgId")

	if queryData != "" {
		mySql := mysql.MysqlDB()
		defer mySql.Close()
		stmt, err := mySql.Query("SELECT freid,flatno FROM fre WHERE ssgId = ?", queryData)
		if err != nil {
			log.Println(err)
			handler.ErrorHandler(c, http.StatusBadRequest, "Query Failed", err)
		}
		defer stmt.Close()

		var Fres []FRE

		for stmt.Next() {
			var result FRE
			err = stmt.Scan(&result.FREId, &result.FlatNo)
			if err != nil {
				handler.ErrorHandler(c, http.StatusInternalServerError, "Query Failed", err)
			}
			Fres = append(Fres, result)
		}
		handler.SuccessHandler(c, http.StatusOK, "success", Fres)

	} else {
		handler.ErrorHandler(c, http.StatusBadRequest, "Query Failed", fmt.Errorf(""))
	}
	// write else block and here
}
