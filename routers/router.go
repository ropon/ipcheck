package routers

import (
	"github.com/gin-gonic/gin"
	v1 "ipinfo/routers/api/v1"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	apiV1 := router.Group("api/v1")
	{
		apiV1.GET("ipinfo/:ipStr", v1.IpInfo)
		apiV1.GET("ipinfo/", v1.IpInfo)
		apiV1.POST("ipinfo", v1.IpInfo)
	}
	return router
}
