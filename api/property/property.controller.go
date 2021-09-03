package property

import (
	"fmt"
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
	print("[file] saved")
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
