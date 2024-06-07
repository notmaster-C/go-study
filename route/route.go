package route

import (
	"github.com/gin-gonic/gin"
	"go-study/file"
)

func InitRoute() {
	r := gin.Default()
	r.Any("/test", handleTest5)
	api := r.Group("/api")
	{
		//api.Any("/test", handleTest5)
		//api.Any("/update", handleSetSswwg)
		//api.Any("/get", handleGet)
		file.InitRoute(api)
	}
	err := r.Run(":8080")
	if err != nil {
		return
	}
}
