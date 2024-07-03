package route

import (
	"go-study/file"
	"go-study/unioffice"

	"github.com/gin-gonic/gin"
)

func InitRoute() {
	r := gin.Default()
	r.Any("/test", handleTest5)
	api := r.Group("/api")
	{
		api.Any("/test", unioffice.TestUnioffice)
		//api.Any("/update", handleSetSswwg)
		//api.Any("/get", handleGet)
		file.InitRoute(api)
	}
	err := r.Run(":8080")
	if err != nil {
		return
	}
}
