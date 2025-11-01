package config

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func NewRedis(viper *viper.Viper) *redis.Client {
	host := viper.GetString("redis.host")
	database := viper.GetInt("redis.database")

	client := redis.NewClient(&redis.Options{
		Addr:     host,
		DB:       database,
		Protocol: 2,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}

	return client
}
