package service

import (
	"order_service/pkg/repository"

	"github.com/redis/go-redis/v9"
)

type OrderService struct{
	repos repository.Order
	redisDB     *redis.Client
}

func NewOrderService(repos repository.Order, redisDB  *redis.Client) *OrderService{
	return &OrderService{
		repos: repos,
		redisDB: redisDB,
	}
}