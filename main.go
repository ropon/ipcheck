package main

import (
	"ipinfo/models"
	"ipinfo/routers"
)

func main() {
	router := routers.InitRouter()
	//初始化数据库
	models.Migrate()
	_ = router.Run(":8099")
	//_ = router.Run("127.0.0.1:8099")
}
