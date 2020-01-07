package v1

import (
	"github.com/gin-gonic/gin"
	"ipcheck/models"
	"ipcheck/utils/tools"
)

func Ping(c *gin.Context) {
	res := models.NewDefaultResult()
	rType := c.DefaultQuery("type", "text")
	defaultIP := tools.GetDefaultIp(c)
	if defaultIP != tools.DefaultIP {
		SecretKey := c.Query("SecretKey")
		//SecretKey无效
		if SecretKey != tools.SecretKey {
			res.ErrCode = 4002
			res.ErrMsg = tools.CodeType[res.ErrCode]
			tools.GetResType(rType, &res, c)
			return
		}
	}
	//可以是域名或IPV4或IPV6
	pingFlag := c.Query("pingFlag")
	err := tools.CheckArg(pingFlag)
	//参数不完整
	if err != nil {
		res.ErrCode = 4001
		res.ErrMsg = tools.CodeType[res.ErrCode]
		tools.GetResType(rType, &res, c)
		return
	}
	if ! tools.CheckIpV4(pingFlag) && !tools.CheckIpV6(pingFlag) && !tools.CheckDomain(pingFlag) {
		res.ErrCode = 4003
		res.ErrMsg = tools.CodeType[res.ErrCode]
		tools.GetResType(rType, &res, c)
		return
	}
	pType := c.DefaultQuery("pType", "ping6")
	//不是ping6就是ping
	if pType != "ping6" {
		pType = "ping"
	}
	//默认读缓存
	uType := c.DefaultQuery("cache", "yes")
	count := c.DefaultQuery("count", "4")
	val, err := tools.RedisGet(pType + pingFlag + count)
	if err != nil || uType != "yes" {
		pingRes, _ := tools.ExecCommand(pType, []string{"-c", count, pingFlag})
		_ = tools.RedisSet(pType+pingFlag+count, pingRes, 60)
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
