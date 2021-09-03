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
	var missingColumn string
	for rows.Next() {
		row := rows.Columns()
		for i, val := range row {
			if headers[i] == val {
				count++
			} else {
				missingColumn = headers[i] //storing missing columns
				break
			}
		}
		break
	}
	var isValidExcel = false
	if count == 7 {
		isValidExcel = true
	} else {

		c.JSON(401, gin.H{
			"message":       "error in excel file",
			"isvalidfile":   isValidExcel,
			"missingColumn": missingColumn,
		})

	}

	mySql := mysql.MysqlDB()

	for i, row := range xlsx.GetRows("All") {

		reversePrice, errReversePrice := strconv.Atoi(row[5])
		emd, errEMD := strconv.Atoi(row[6])

		if errReversePrice == nil || errEMD == nil { //this should be and

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
