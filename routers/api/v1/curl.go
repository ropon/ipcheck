package v1

import (
	"github.com/gin-gonic/gin"
	"ipcheck/models"
	"ipcheck/utils/tools"
)

func Curl(c *gin.Context) {
	res := models.NewDefaultResult()
	rType := c.DefaultQuery("type", "text")
	cType := c.DefaultQuery("ctype", "-6")
	cStatus := c.Query("cstatus")
	//可以是域名或IPV4
	curlFlag := c.Query("curlflag")
	//默认读缓存
	uType := c.DefaultQuery("cache", "yes")
	err := tools.CheckArg(curlFlag)
	//参数不完整
	if err != nil {
		res.ErrCode = 4001
		res.ErrMsg = tools.CodeType[res.ErrCode]
		tools.GetResType(rType, &res, c)
		return
	}
	val, err := tools.RedisGet("curl:" + curlFlag)
	if err != nil || uType != "yes" {
		curlRes, _ := tools.ExecCommand("curl", []string{cType, cStatus, curlFlag})
		_ = tools.RedisSet("curl:"+curlFlag, curlRes, 60)
		res.Data = curlRes
	} else {
		res.Data = val
	}
	if rType == "json" {
		c.JSON(200, res)
	} else {
		c.String(200, res.Data.(string))
	}
	return
}
