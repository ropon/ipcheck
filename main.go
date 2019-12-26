package main

import "ipcheck/routers"

func main() {
	router := routers.InitRouter()
	//初始化数据库
	//测试修改
	_ = router.Run(":8090")
	//_ = router.Run("127.0.0.1:8090")
}
