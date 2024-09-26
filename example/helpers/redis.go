package helpers

import (
	"github.com/tiant-developer/golib/example/conf"
	"github.com/tiant-developer/golib/redis"
)

// 推荐，直接使用
var RedisClient *redis.Redis

// 初始化redis
func InitRedis() {
	c := conf.GlobalC.Redis["default"]
	var err error
	RedisClient, err = redis.InitRedisClient(c)
	if err != nil || RedisClient == nil {
		panic("init redis failed!")
	}
}

func CloseRedis() {
	_ = RedisClient.Close()
}
