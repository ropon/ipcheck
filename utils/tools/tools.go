package tools

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"ipcheck/models"
	"os/exec"
	"time"
)

//定义错误代码说明
var CodeType = map[int]string{
	4001: "参数不完整",
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

func ExecCommand(commandName string, params []string, ) (res string, err error) {
	stdout, err := exec.Command(commandName, params...).Output()
	if err != nil {
		return
	}
	res = string(stdout)
	return
}
