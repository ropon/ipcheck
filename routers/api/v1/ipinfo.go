package v1

import (
	"encoding/binary"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"ipinfo/models"
	"ipinfo/utils/tools"
	"net"
)

func IpInfo(c *gin.Context) {
	res := models.NewDefaultResult()
	if c.Request.Method == "POST" {
		tools.GetIpData()
	} else if c.Request.Method == "GET" {
		ipStr := c.Param("ipStr")
		rType := c.DefaultQuery("type", "text")
		if ipStr == "" {
			ipStr = tools.GetDefaultIp(c)
		}
		ip := net.ParseIP(ipStr)
		//ip格式不正确
		if ip == nil || ip.To4() == nil {
			res.ErrCode = 4003
			res.ErrMsg = tools.CodeType[res.ErrCode]
			tools.GetResType(rType, &res, c)
			return
		}
		//将ip转换uint32
		ipNum := binary.BigEndian.Uint32([]byte(ip.To4()))
		tmpIpInfo := models.IpInfo{}
		val, err := tools.RedisGet(ipStr)
		if err != nil {
			tmpIpInfo = models.IpInfoGet(ipNum)
			if tmpIpInfo.Zone == "" {
				res.ErrCode = 4004
				res.ErrMsg = tools.CodeType[res.ErrCode]
				tools.GetResType(rType, &res, c)
				return
			}
			ipInfoByte, _ := json.Marshal(tmpIpInfo)
			_ = tools.RedisSet(ipStr, string(ipInfoByte), 120)
		} else {
			_ = json.Unmarshal([]byte(val.(string)), &tmpIpInfo)
		}
		tmpIpInfo.Ip = ipStr
		res.Data = tmpIpInfo
		if rType == "json" {
			c.JSON(200, res)
		} else {
			c.String(200, ipStr+" "+tmpIpInfo.Zone+" "+tmpIpInfo.Isp)
		}
		return
	}
	c.JSON(200, res)
}
