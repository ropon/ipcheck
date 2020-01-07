package routers

import (
	"github.com/gin-gonic/gin"
	"io"
	v1 "ipcheck/routers/api/v1"
	"os"
)

func InitRouter() *gin.Engine {
	gin.DisableConsoleColor()
	f, _ := os.Create("./ipcheck.log")
	gin.DefaultWriter = io.MultiWriter(f)
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	apiV1 := router.Group("api/v1")
	{
		apiV1.GET("ping", v1.Ping)
		apiV1.GET("curl", v1.Curl)
	}
	return router
}
