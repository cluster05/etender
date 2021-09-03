package property

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

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

	var propertyDTO []PropertyDTO

	/*
		change sheet title name
	*/
	rows, err := xlsx.Rows("All")
	if err != nil {
		log.Fatal(err)
	}

	var count = 0
	headers := []string{"Division", "Station", "Sector", "Group", "Flat No.", "Reserve price", "EMD"}
	for rows.Next() {
		row := rows.Columns()
		var header_slice = row
		log.Println(header_slice)
		for _, val := range row {
			if contains(headers, val) {
				count++
			}
		}
		break
	}
	if count == 7 {
		log.Println("Valid xlsx")
	} else {
		log.Println("Invalid xlsx")
	}

	for _, row := range xlsx.GetRows("Sheet1") {

		_ = row
		var temp PropertyDTO
		/*
			1. you can get col by row[i]
			2. validate data before insert
		*/

		propertyDTO = append(propertyDTO, temp)

	}

}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
