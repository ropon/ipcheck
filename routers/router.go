package routers

import (
	"github.com/gin-gonic/gin"
	v1 "ipcheck/routers/api/v1"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	apiV1 := router.Group("api/v1")
	{
		apiV1.GET("ping", v1.Ping)
		apiV1.GET("curl", v1.Curl)
	}
	return router
}
