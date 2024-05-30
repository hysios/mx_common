package common

import (
	"github.com/creasty/defaults"
	"github.com/hysios/mx/config"
	"github.com/redis/go-redis/v9"
)

type RedisOption struct {
	Network  string `default:"tcp"`
	Addr     string `default:"localhost:6379"`
	DB       int
	Password string
}

func GetRedis(cfg *config.Config) RedisOption {
	opts := RedisOption{
		Network:  cfg.Str("redis.network"),
		Addr:     cfg.Str("redis.addr"),
		DB:       cfg.Int("redis.db"),
		Password: cfg.Str("redis.password"),
	}

	defaults.Set(&opts)
	return opts
}

func OpenRedis(cfg *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Network: cfg.Str("redis.network"),
		Addr:    cfg.Str("redis.addr"),
		DB:      cfg.Int("redis.db"),
		// Username: cfg.Str("redis.username"),
		Password: cfg.Str("redis.password"),
	})
}
