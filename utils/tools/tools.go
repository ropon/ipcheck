package tools

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"ipcheck/models"
	"net"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

//定义错误代码说明
var CodeType = map[int]string{
	4001: "参数不完整",
	4002: "SecretKey无效",
	4003: "IP或域名格式不正确",
}

var (
	redisDb   *redis.Client
	SecretKey = "*********"
	DefaultIP = "127.0.0.1"
	reg       = "^((xn--)?[A-Za-z0-9*]{1,100}\\.){1,8}((xn--)?[A-Za-z0-9]){1,24}$"
)

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

func CheckIpV4(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil || ip.To4() == nil {
		return false
	}
	return true
}

func CheckIpV6(ipStr string) bool {
	return strings.Count(ipStr, ":") >= 2
}

func CheckDomain(domainStr string) bool {
	match, _ := regexp.MatchString(reg, domainStr)
	return match
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

func GetResType(rType string, res *models.Result, c *gin.Context) {
	if rType == "json" {
		c.JSON(200, *res)
	} else {
		c.String(200, res.ErrMsg)
	}
}

func ExecCommand(commandName string, params []string, ) (res string, err error) {
	//返回标准输出和错误
	out, err := exec.Command(commandName, params...).CombinedOutput()
	if err != nil {
		return
	}
	res = string(out)
	return
}
