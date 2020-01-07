package v1

import (
	"github.com/gin-gonic/gin"
	"ipcheck/models"
	"ipcheck/utils/tools"
)

func Curl(c *gin.Context) {
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
	//可以是域名或IPV4
	curlFlag := c.Query("curlFlag")
	err := tools.CheckArg(curlFlag)
	//参数不完整
	if err != nil {
		res.ErrCode = 4001
		res.ErrMsg = tools.CodeType[res.ErrCode]
		tools.GetResType(rType, &res, c)
		return
	}
	if ! tools.CheckIpV4(curlFlag) && !tools.CheckDomain(curlFlag) {
		res.ErrCode = 4003
		res.ErrMsg = tools.CodeType[res.ErrCode]
		tools.GetResType(rType, &res, c)
		return
	}
	cType := c.DefaultQuery("cType", "-6")
	cStatus := c.Query("cStatus")
	//默认读缓存
	uType := c.DefaultQuery("cache", "yes")
	val, err := tools.RedisGet(cType + curlFlag + cStatus)
	if err != nil || uType != "yes" {
		curlRes, _ := tools.ExecCommand("curl", []string{cType, cStatus, curlFlag})
		_ = tools.RedisSet(cType+curlFlag+cStatus, curlRes, 60)
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
