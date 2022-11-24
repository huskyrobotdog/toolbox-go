package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/huskyrobotdog/toolbox-go/inner"
)

var Instance *redis.Client

type Config struct {
	Address  string `json:"address"`
	Password string `json:"password"`
	Database int    `json:"database"`
}

func Initialization(config *Config) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Address,
		Password: config.Password,
		DB:       config.Database,
	})
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		inner.Fatal(err.Error())
	}
	Instance = client
	inner.Debug("cache initialization complete")
}
