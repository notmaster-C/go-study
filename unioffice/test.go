package unioffice

import (
	"errors"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unidoc/unioffice/spreadsheet"
)

func TestUnioffice(c *gin.Context) {
	var (
		err    error
		result = gin.H{}
		wb     *spreadsheet.Workbook
		f      multipart.File
	)
	defer func() {

		if err != nil {
			result["error"] = err
		}
		c.JSON(http.StatusOK, result)
	}()
	file, err := c.FormFile("file")
	if err != nil {
		return
	}
	f, err = file.Open()
	defer f.Close()

	if err != nil {
		return
	}
	wb, err = spreadsheet.Read(f, file.Size)
	defer wb.Close()

	if err != nil {
		return
	}
	// ExcelData := make(map[string]string)
	for _, v := range wb.Sheets() {
		if len(v.Rows()) < 3 {
			err = errors.New("内容格式不满足要求")
			return
		}
		// rows2Json(v.Rows())
		result["content"] = v.Rows()
	}
}
func rows2Jso1n(rows []spreadsheet.Row) {
	// headers := rows[0]
	// fieldsName := rows[1]
	// jsonData := make([]map[string]interface{}, len(rows)-2)
	// for _,row:=rows[2:]{
	// 	fmt.Println(row)
	// }
}
