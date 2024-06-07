package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 输出结果
func OutputResult(c *gin.Context, err error, result gin.H) {
	if result == nil {
		result = gin.H{}
	}

	if err != nil {
		fmt.Errorf("bad request for %s with :%v", c.Request.URL.Path, err)
		result["code"] = -1
		result["message"] = err
		result["success"] = false
	} else {
		result["code"] = 200
		result["success"] = true
	}

	c.JSON(http.StatusOK, result)
}
