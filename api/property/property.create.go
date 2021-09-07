package property

import (
	"etender/api/handler"
	"etender/mysql"
	"fmt"
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
		handler.ErrorHandler(c, http.StatusBadRequest, "error in saving excel file", err)
		return
	}

	extension := filepath.Ext(file.Filename)

	newFileName := uuid.New().String() + extension
	if err := c.SaveUploadedFile(file, "./excelfiles/"+newFileName); err != nil {
		handler.ErrorHandler(c, http.StatusBadRequest, "error in saving excel file", err)
		return
	}
	saveExcelInDB(newFileName, c)
}

func saveExcelInDB(filename string, c *gin.Context) {

	xlsx, err := excelize.OpenFile("./excelfiles/" + filename)

	if err != nil {
		handler.ErrorHandler(c, http.StatusBadRequest, "error reading excel file", err)
		return
	}

	rows, err := xlsx.Rows("All")
	if err != nil {
		handler.ErrorHandler(c, http.StatusBadRequest, "unable to read excel file", err)
		return
	}

	var count = 0
	headers := []string{"Division", "Station", "Sector", "Group", "Flat No.", "Reserve price", "EMD"}

	var checkForMissingColumns string = " Following columns are missing : "

	for rows.Next() {
		row := rows.Columns()
		for i, val := range row {
			if headers[i] == val {
				count++
			} else {
				checkForMissingColumns += headers[i] + " "
			}
		}
		break
	}

	if count != 7 {
		errorString := "look like you tempered with main template that provided by us." + checkForMissingColumns
		handler.ErrorHandler(c, http.StatusBadRequest, errorString, err)
		return
	}

	mySql := mysql.MysqlDB()
	defer mySql.Close()

	stmtDivision, errDivison := mySql.Prepare("INSERT INTO division(name) VALUES(?)")
	stmtSSG, errSSG := mySql.Prepare("INSERT INTO ssg(station,sector,pgroup,reserveprice,emd,uniquestream,divisionId) VALUES(?,?,?,?,?,?,?)")
	stmtFRE, errFRE := mySql.Prepare("INSERT INTO fre(flatno,ssgId,uniquefre) VALUES(?,?,?)")
	defer stmtDivision.Close()
	defer stmtSSG.Close()
	defer stmtFRE.Close()

	if errDivison != nil && errSSG != nil && errFRE != nil {
		handler.ErrorHandler(c, http.StatusInternalServerError, "error in query. contact developer :).", fmt.Errorf(":("))
		defer mySql.Close()
		return
	}

	mapperDivison := make(map[string]int)
	mapperSSG := make(map[string]int)

	counter := struct {
		skipedEntry   int
		insertedEntry int
		totalEntry    int
	}{}

	for i, row := range xlsx.GetRows("All") {

		reversePrice, errReversePrice := strconv.Atoi(row[5])
		emd, errEMD := strconv.Atoi(row[6])

		if errReversePrice == nil && errEMD == nil {

			var temp = PropertyDTO{
				Division:     strings.ToLower(row[0]),
				Station:      strings.ToLower(row[1]),
				Sector:       row[2],
				Group:        row[3],
				FlatNo:       row[4],
				ReversePrice: reversePrice,
				EMD:          emd,
			}

			if isValidEntry(temp) {

				uniquestream := temp.Station + "<>" + temp.Sector + "<>" + temp.Group

				if mapperDivison[temp.Division] == 0 {

					resDivision, errExecDivision := stmtDivision.Exec(temp.Division)

					if errExecDivision == nil {
						lastInsetedIdDivision, err := resDivision.LastInsertId()

						if err != nil {
							fmt.Printf("[PROPERTY][CREATE][logger] [%v] Error %v  \n", i, err.Error())
						}

						mapperDivison[temp.Division] = int(lastInsetedIdDivision)
					} else {
						var tempDivisionId int
						err := mySql.QueryRow("SELECT divisionId from division WHERE name= ? ", temp.Division).Scan(&tempDivisionId)

						if err != nil {
							if err != nil {
								fmt.Printf("[PROPERTY][CREATE][logger] [%v] Error %v  \n", i, err.Error())
							}
						}
						mapperDivison[temp.Division] = tempDivisionId
					}
				}

				if mapperSSG[uniquestream] == 0 {

					resSSG, errExceSSG := stmtSSG.Exec(temp.Station, temp.Sector, temp.Group, temp.ReversePrice, temp.EMD, uniquestream, mapperDivison[temp.Division])

					if errExceSSG == nil {

						lastInsetedIdDivision, err := resSSG.LastInsertId()

						if err != nil {
							if err != nil {
								fmt.Printf("[PROPERTY][CREATE][logger] [%v] Error %v  \n", i, err.Error())
							}
						}
						mapperSSG[uniquestream] = int(lastInsetedIdDivision)

					} else {

						var tempSSGId int
						err := mySql.QueryRow("SELECT ssgId from ssg WHERE uniquestream= ? ", uniquestream).Scan(&tempSSGId)

						if err != nil {
							if err != nil {
								fmt.Printf("[PROPERTY][CREATE][logger] [%v] Error %v  \n", i, err.Error())
							}
						}
						mapperSSG[uniquestream] = tempSSGId
					}

				}

				uniquefre := temp.FlatNo + "<>" + strconv.Itoa(mapperSSG[uniquestream])

				_, errExecFRE := stmtFRE.Exec(temp.FlatNo, mapperSSG[uniquestream], uniquefre)

				if errExecFRE != nil {
					counter.skipedEntry++
					counter.totalEntry++
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

	response := make(map[string]int)

	response["skippedEntry"] = counter.skipedEntry - 1
	response["totalEntry"] = counter.totalEntry - 1
	response["insertedEntry"] = counter.insertedEntry

	handler.SuccessHandler(c, http.StatusCreated, "entry inseted", response)

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
