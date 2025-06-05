package service

import (
	"context"
	"encoding/json"
	"log"
	productservice "product_service"
	"product_service/pkg/repository"

	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
)


type ProductService struct{
	repos repository.Product
	minIO *minio.Client
	redisDB *redis.Client
}

func MewProductService(repos repository.Product, minIO *minio.Client, redisDB *redis.Client) *ProductService{
	return &ProductService{
		repos: repos,
		minIO: minIO,
		redisDB: redisDB,
	}
}

func (s *ProductService) CatregoList()([]productservice.CategorySerializer, error){
	data_r, err := s.redisDB.Get(context.Background(), RedisCategory).Result()
	if err != nil{
		log.Printf("cache error %s ", err.Error())
		return s.repos.CatregoListPostgres(true)
	}
	if data_r == "" {
		log.Println("category cache empty")
		return s.repos.CatregoListPostgres(true)
	}
	var data []productservice.CategorySerializer

	err = json.Unmarshal([]byte(data_r), &data)

	if err != nil{
		return data, err 
	}

	return data, nil
	
}