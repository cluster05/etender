package property

import (
	"database/sql"
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

/*
	TODO :
	handler error message in one place
	handler success message in one place
*/

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
			check for error in convertion
			if error skip
			else good to process
		*/
		if errReversePrice == nil && errEMD == nil {

			var temp = PropertyDTO{
				Division:     strings.ToLower(row[0]),
				Station:      row[1],
				Sector:       row[2],
				Group:        row[3],
				FlatNo:       row[4],
				ReversePrice: reversePrice,
				EMD:          emd,
			}

			if isValidEntry(temp) {

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

				var resDivision sql.Result
				var errExecDivision error

				if mapperDivison[temp.Division] == 0 {

					fmt.Println("Insertion avoided")
					resDivision, errExecDivision = stmtDivision.Exec(temp.Division)

					if errExecDivision == nil {
						lastInsetedIdDivision, err := resDivision.LastInsertId()

						// implement skip query
						if err != nil {
							fmt.Printf("[Testsql] [%v] Error %v  \n", i, err.Error())
						}

						mapperDivison[temp.Division] = int(lastInsetedIdDivision)
					} else {
						var tempDivisionId int
						err := mySql.QueryRow("SELECT divisionId from division WHERE name= ? ", temp.Division).Scan(&tempDivisionId)
						if err != nil {
							fmt.Printf("[Testsql] [%v] Error %v  \n", i, err.Error())
						}

						fmt.Printf("[Duplicate Entry Found Make Map for it] %v\n", tempDivisionId)
						mapperDivison[temp.Division] = tempDivisionId
					}
				}

				_, errExceSSG := stmtSSG.Exec(temp.Station, temp.Sector, temp.Group, uniquestream, mapperDivison[temp.Division])
				_, errExecFRE := stmtFRE.Exec(temp.FlatNo, temp.ReversePrice, temp.EMD, 1)

				// implement skip query
				if errExecDivision != nil && errExceSSG != nil && errExecFRE != nil {
					//implement skip entry if data is repeated
					fmt.Printf("[Testsql] [%v] Insertion \nDivision Error %v\n SSG Error %v\n FRE Error %v\n", i, errExecDivision.Error(), errExceSSG.Error(), errExecFRE.Error())
				}

			} else {
				// implement skip entry

			}
		} else {
			// implement skip entry
		}

		// implement total entry
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
