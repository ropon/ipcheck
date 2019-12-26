package v1

import (
	"github.com/gin-gonic/gin"
	"ipcheck/models"
	"ipcheck/utils/tools"
)

func Ping(c *gin.Context) {
	res := models.NewDefaultResult()
	rType := c.DefaultQuery("type", "text")
	pType := c.DefaultQuery("ptype", "ping6")
	//可以是域名或IPV4或IPV6
	pingFlag := c.Query("pingflag")
	//默认读缓存
	uType := c.DefaultQuery("cache", "yes")
	err := tools.CheckArg(pingFlag)
	//参数不完整
	if err != nil {
		res.ErrCode = 4001
		res.ErrMsg = tools.CodeType[res.ErrCode]
		tools.GetResType(rType, &res, c)
		return
	}
	val, err := tools.RedisGet(pType + pingFlag)
	if err != nil || uType != "yes" {
		pingRes, _ := tools.ExecCommand(pType, []string{"-c", "4", pingFlag})
		_ = tools.RedisSet(pType+pingFlag, pingRes, 60)
		res.Data = pingRes
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
