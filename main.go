package main

import (
	"ipcheck/routers"
)

func main() {
	router := routers.InitRouter()
	_ = router.Run(":8090")
}
