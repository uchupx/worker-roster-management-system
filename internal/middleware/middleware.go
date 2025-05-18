package middleware

import redis "github.com/redis/go-redis/v9"

type Middleware struct {
	Redis *redis.Client
}

type Config struct {
	Redis *redis.Client
}

func New(conf Config) *Middleware {
	return &Middleware{
		Redis: conf.Redis,
	}
}
