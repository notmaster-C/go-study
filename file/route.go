package file

import (
	"github.com/gin-gonic/gin"
)

var privateMetaData = map[string]string{"x-amz-acl": "private"}

func InitRoute(api *gin.RouterGroup) {
	api.Any("/file/upload", handleUpload)
	api.Group("/file")
	{
		api.Any("/upload", handleUpload)
	}
}
