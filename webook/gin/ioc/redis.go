package ioc

import (
	"github.com/redis/go-redis/v9"
	"new_home/webook/gin/config"
)

func InitRedis() redis.Cmdable {
	redisclient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	return redisclient
}
