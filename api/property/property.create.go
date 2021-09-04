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

	mapperDivison := make(map[string]int)
	mapperSSG := make(map[string]int)

	type Counter struct {
		skipedEntry   int
		insertedEntry int
		totalEntry    int
	}
	var counter Counter

	for i, row := range xlsx.GetRows("All") {

		reversePrice, errReversePrice := strconv.Atoi(row[5])
		emd, errEMD := strconv.Atoi(row[6])

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

				var errExecDivision error
				if mapperDivison[temp.Division] == 0 {

					resDivision, errExecDivision := stmtDivision.Exec(temp.Division)

					if errExecDivision == nil {
						lastInsetedIdDivision, err := resDivision.LastInsertId()

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
						mapperDivison[temp.Division] = tempDivisionId
					}
				}

				var errExceSSG error
				if mapperSSG[uniquestream] == 0 {

					resSSG, errExceSSG := stmtSSG.Exec(temp.Station, temp.Sector, temp.Group, uniquestream, mapperDivison[temp.Division])

					if errExceSSG == nil {

						lastInsetedIdDivision, err := resSSG.LastInsertId()

						if err != nil {
							fmt.Printf("[Testsql] [%v] Error %v  \n", i, err.Error())
						}
						mapperSSG[uniquestream] = int(lastInsetedIdDivision)

					} else {

						var tempSSGId int
						err := mySql.QueryRow("SELECT ssgId from ssg WHERE uniquestream= ? ", uniquestream).Scan(&tempSSGId)

						if err != nil {
							fmt.Printf("[Testsql] [%v] Error %v  \n", i, err.Error())
						}
						mapperSSG[uniquestream] = tempSSGId
					}

				}

				_, errExecFRE := stmtFRE.Exec(temp.FlatNo, temp.ReversePrice, temp.EMD, mapperSSG[uniquestream])

				if errExecDivision != nil && errExceSSG != nil && errExecFRE != nil {
					counter.skipedEntry++
					fmt.Printf("[Testsql] [%v] Insertion \nDivision Error %v\n SSG Error %v\n FRE Error %v\n", i, errExecDivision.Error(), errExceSSG.Error(), errExecFRE.Error())
					continue
				}

				counter.insertedEntry++

			} else {
				counter.skipedEntry++
			}
		} else {
			counter.skipedEntry++
		}

		fmt.Printf("[Test SQL] Inserting value index [%v]\n", i)
		counter.totalEntry++

	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "entery inserted",
		"skippedEntry":  counter.skipedEntry,
		"totalEntry":    counter.totalEntry,
		"insertedEntry": counter.insertedEntry,
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
