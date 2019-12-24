package tools

import (
	"encoding/binary"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"io/ioutil"
	"ipinfo/models"
	"net"
	"strings"
	"time"
)

//加密key
var SecretKey = "West.cn209129"

//定义错误代码说明
var CodeType = map[int]string{
	4001: "参数不完整",
	4002: "SecretKey无效",
	4003: "IP格式不正确",
	4004: "记录不存在",
}

type OpenIdRes struct {
	SessionKey string `json:"session_key"`
	OpenId     string `json:"openid"`
}

var redisDb *redis.Client

func init() {
	redisDb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	_, _ = redisDb.Ping().Result()
}

func RedisSet(key string, val interface{}, expiration int64) (err error) {
	err = redisDb.Set(key, val, time.Second*time.Duration(expiration)).Err()
	if err != nil {
		return
	}
	return
}

func RedisGet(key string) (val interface{}, err error) {
	val, err = redisDb.Get(key).Result()
	if err == redis.Nil {
		err = fmt.Errorf("键%s对应值不存在", key)
	} else if err != nil {
		return
	}
	return
}

//检查参数
func CheckArg(args ...string) (err error) {
	aLen := len(args)
	for i := 0; i < aLen; i++ {
		if args[i] == "" {
			err = fmt.Errorf(CodeType[4001])
			return
		}
	}
	return
}

func GetResType(rType string, res *models.Result, c *gin.Context) {
	if rType == "json" {
		c.JSON(200, *res)
	} else {
		c.String(200, res.ErrMsg)
	}
}

func GetDefaultIp(c *gin.Context) string {
	remoteAddr := c.Request.RemoteAddr
	if ip := c.Request.Header.Get("HTTP_X_FORWARDED_FOR"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}
	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}
	return remoteAddr
}

func convertUint32(ipStr string) (ipNum uint32, err error) {
	ip := net.ParseIP(ipStr)
	//ip格式不正确
	if ip == nil || ip.To4() == nil {
		err = fmt.Errorf("ip格式不正确")
		return
	}
	ipNum = binary.BigEndian.Uint32([]byte(ip.To4()))
	return
}

func GetIpData() {
	data, err := ioutil.ReadFile("./ipdata.txt")
	if err != nil {
		err = fmt.Errorf("打开配置文件%s失败", "./ipdata.txt")
		return
	}
	lineSlice := strings.Split(string(data), "\r\n")
	for _, v := range lineSlice[500001:] {
		tmpIpInfo := models.IpInfo{}
		vSlice := strings.Split(v, " ")
		var newSlice []string
		for _, v2 := range vSlice {
			if v2 != "" {
				newSlice = append(newSlice, v2)
			}
		}
		tmpIpInfo.Start, err = convertUint32(newSlice[0])
		if err != nil {
			break
		}
		tmpIpInfo.End, err = convertUint32(newSlice[1])
		if err != nil {
			break
		}
		tmpIpInfo.Zone = newSlice[2]
		if len(newSlice) > 3 {
			tmpIpInfo.Isp = newSlice[3]
		} else {
			tmpIpInfo.Isp = ""
		}
		tmpIpInfo.Insert()
	}
}
