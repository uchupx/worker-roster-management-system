package database

import (
	"github.com/redis/go-redis/v9"
	"github.com/uchupx/worker-roster-management-system/config"
)


var redisClient *redis.Client

func GetRedisClient(conf config.Config) *redis.Client {
	if redisClient != nil {
		return redisClient
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Host + ":" + conf.Redis.Port,
		Password: conf.Redis.Password,
		DB:       0,
	})

	return redisClient
}
