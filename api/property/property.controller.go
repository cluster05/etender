package property

import (
	"etender/mysql"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetSsg(c *gin.Context) {
	query_data := c.Param("division")
	fmt.Println(query_data)
	if query_data != "" {
		fmt.Printf("%v", query_data)
		mySql := mysql.MysqlDB()
		stmt, err := mySql.Query("Select ssgid,station,sector,pgroup from ssg where divisionId = ?", query_data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Query Failed",
				"err":     err.Error(),
			})
		}
		defer stmt.Close()
		var Ssgs []Ssg
		for stmt.Next() {
			var result Ssg
			if err := stmt.Scan(&result.SsgId, &result.Station, &result.Sector,
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
func GetFre(c *gin.Context) {
	query_data := c.Param("ssgid")
	if query_data != "" {
		mySql := mysql.MysqlDB()
		stmt, err := mySql.Query("Select freid,flatno,reserveprice,emd from fre where ssgId = ?", query_data)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Query Failed",
				"err":     err.Error(),
			})
		}
		defer stmt.Close()
		var Fres []Fre
		for stmt.Next() {
			var result Fre
			err = stmt.Scan(&result.FreId, &result.FlatNo, &result.ReservePrice, &result.EMD)
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
		defer mySql.Close()
	}
}
func GetDivision(c *gin.Context) {
	mySql := mysql.MysqlDB()
	stmt, err := mySql.Query("Select divisionId,name from division")
	var Divisions []Division
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Query Failed",
			"err":     err.Error(),
		})
	}
	defer stmt.Close()
	for stmt.Next() {
		var result Division
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

func SaveExcel(c *gin.Context) {

	file, err := c.FormFile("file")

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No file is received",
		})
		return
	}

	extension := filepath.Ext(file.Filename)

	newFileName := uuid.New().String() + extension
	if err := c.SaveUploadedFile(file, "./excelfiles/"+newFileName); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to save the file",
			"error":   err.Error(),
		})
		return
	}
	log.Println("[file] saved")
	saveExcelInDB(newFileName, c)
}

func saveExcelInDB(filename string, c *gin.Context) {

	xlsx, err := excelize.OpenFile("./excelfiles/" + filename)
	if err != nil {
		fmt.Println(err.Error())
		print("failed to read")
		return
	}

	// var propertyDTO []PropertyDTO

	rows, err := xlsx.Rows("All")
	if err != nil {
		fmt.Println(err.Error())
	}

	var count = 0
	headers := []string{"Division", "Station", "Sector", "Group", "Flat No.", "Reserve price", "EMD"}
	for rows.Next() {
		row := rows.Columns()
		for i, val := range row {
			if headers[i] == val {
				count++
			}
		}
		break
	}
	var isValidExcel = false
	if count == 7 {
		isValidExcel = true
	} else {

		c.JSON(401, gin.H{
			"message":     "error in excel file",
			"isvalidfile": isValidExcel,
		})

	}

	mySql := mysql.MysqlDB()

	for i, row := range xlsx.GetRows("All") {

		reversePrice, errReversePrice := strconv.Atoi(row[5])
		emd, errEMD := strconv.Atoi(row[6])

		if errReversePrice == nil || errEMD == nil {

			var temp = PropertyDTO{
				Division:     row[0],
				Station:      row[1],
				Sector:       row[2],
				Group:        row[3],
				FlatNo:       row[4],
				ReversePrice: reversePrice,
				EMD:          emd,
			}

			stmt, err := mySql.Prepare("INSERT INTO division(name) VALUES(?)")
			if err != nil {
				fmt.Printf("[Testsql] Error %v\n", err.Error())
			}
			defer stmt.Close()

			result, err := stmt.Exec(temp.Division)
			fmt.Printf("[Testsql] [%v] Insertion log %v\n", i, result)
			if err != nil {
				// fmt.Printf("[Testsql] [%v] Insertion Error %v\n", i, err.Error())
			}

		}
	}

	c.JSON(201, gin.H{
		"message": "entery inserted",
	})
	defer mySql.Close()

}
