package property

import (
	"etender/mysql"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

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

	type SSGView struct {
		uniquestream string
		ssgId        int
	}

	mapperDivison := make(map[string]int)
	// mapperSSG := make(map[string]SSGView)

	for i, row := range xlsx.GetRows("All") {

		reversePrice, errReversePrice := strconv.Atoi(row[5])
		emd, errEMD := strconv.Atoi(row[6])

		/*
			handler uppercase lower case to avoid duplicate entry
		*/
		var temp = PropertyDTO{
			Division:     strings.ToLower(row[0]),
			Station:      row[1],
			Sector:       row[2],
			Group:        row[3],
			FlatNo:       row[4],
			ReversePrice: reversePrice,
			EMD:          emd,
		}

		if errReversePrice == nil && errEMD == nil && isValidEntry(temp) {

			uniquestream := temp.Station + "<>" + temp.Sector + "<>" + temp.Group

			stmtDivision, errDivison := mySql.Prepare("INSERT INTO division(name) VALUES(?)")
			stmtSSG, errSSG := mySql.Prepare("INSERT INTO ssg(station,sector,pgroup,uniquestream,divisionId) VALUES(?,?,?,?,?)")
			stmtFRE, errFRE := mySql.Prepare("INSERT INTO fre(flatno,reserveprice,emd,ssgId) VALUES(?,?,?,?)")
			defer stmtDivision.Close()
			defer stmtSSG.Close()
			defer stmtFRE.Close()

			if errDivison != nil && errSSG != nil && errFRE != nil {
				c.JSON(201, gin.H{
					"message": "error in query writring contact developer",
				})
				defer mySql.Close()
			}

			resDivison, errExecDivision := stmtDivision.Exec(temp.Division)

			/*
				if not present in map search it first in database
			*/

			if errExecDivision == nil {
				lastInsetedIdDivision, err := resDivison.LastInsertId()

				if err != nil {
					fmt.Printf("[Testsql] [%v] Error %v  \n", i, err.Error())
				}

				mapperDivison[temp.Division] = int(lastInsetedIdDivision)
			}

			fmt.Printf("Testsql] [%v] lasteInsertionId Division is [%v]\n", i, mapperDivison[temp.Division])

			_, errExceSSG := stmtSSG.Exec(temp.Station, temp.Sector, temp.Group, uniquestream, mapperDivison[temp.Division])
			_, errExecFRE := stmtFRE.Exec(temp.FlatNo, temp.ReversePrice, temp.EMD, 1)

			if errExecDivision != nil && errExceSSG != nil && errExecFRE != nil {
				fmt.Printf("[Testsql] [%v] Insertion \nDivision Error %v\n SSG Error %v\n FRE Error %v\n", i, errExecDivision.Error(), errExceSSG.Error(), errExecFRE.Error())
			}

		}
	}

	c.JSON(201, gin.H{
		"message": "entery inserted",
	})
	defer mySql.Close()

}

func isValidEntry(temp PropertyDTO) bool {
	if len(temp.Division) > 0 &&
		len(temp.Station) > 0 &&
		len(temp.Sector) > 0 &&
		len(temp.Group) > 0 &&
		len(temp.FlatNo) > 0 &&
		temp.ReversePrice > 0 &&
		temp.EMD > 0 {
		return true
	}
	return false
}
