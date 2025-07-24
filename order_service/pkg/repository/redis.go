package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func NewRedisConnect(host string, port int, password string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       0,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Printf("failed to connect to redis server: %s\n", err.Error())
		return nil, err
	}
	return client, nil

}
