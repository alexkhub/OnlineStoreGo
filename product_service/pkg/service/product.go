package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	productservice "product_service"
	"product_service/pkg/repository"
	"time"

	"github.com/redis/go-redis/v9"
)

type ProductService struct {
	repos       repository.Product
	redisDB     *redis.Client
	minioClient repository.MinIO
}

func MewProductService(repos repository.Product, redisDB *redis.Client, minioClient repository.MinIO) *ProductService {
	return &ProductService{
		repos:       repos,
		redisDB:     redisDB,
		minioClient: minioClient,
	}
}

func (s *ProductService) CategoryList() ([]productservice.CategorySerializer, error) {
	data_r, err := s.redisDB.Get(context.Background(), RedisCategory).Result()
	if err != nil {
		log.Printf("cache error %s ", err.Error())
		return s.repos.CategoryListPostgres(true)
	}
	if data_r == "" {
		log.Println("category cache empty")
		return s.repos.CategoryListPostgres(true)
	}
	var data []productservice.CategorySerializer
	err = json.Unmarshal([]byte(data_r), &data)

	if err != nil {
		return data, err
	}

	return data, nil

}

func (s *ProductService) ProductList() ([]productservice.ProductListSerailizer, error) {
	cacheData, err := s.redisDB.Get(context.Background(), "products").Result()

	if err == nil {
		var product []productservice.ProductListSerailizer
		if err := json.Unmarshal([]byte(cacheData), &product); err == nil {

			return product, nil
		}
	}
	data, err := s.repos.ProductListPostgres()
	if err != nil {
		return []productservice.ProductListSerailizer{}, err
	}
	image_names := make([]string, 0, len(data))

	for _, product := range data {
		if product.Image.Valid {
			image_names = append(image_names, fmt.Sprintf("product%d/%s", product.Id, product.Image.String))
		}
	}

	urls, err := s.minioClient.GetMany("product", image_names)
	if err != nil {
		return []productservice.ProductListSerailizer{}, err
	}
	for indx, product := range data {
		if product.Image.Valid {
			url, ok := urls[fmt.Sprintf("product%d/%s", product.Id, product.Image.String)]
			if ok {
				data[indx].ImageLink.SetValid(url)
			}
		}
	}
	newCache, _ := json.Marshal(data)
	if err := s.redisDB.Set(context.Background(), "products", newCache, 2*time.Hour).Err(); err != nil {
		fmt.Printf("failed to set data, error: %s", err.Error())
	}

	return data, nil
}

func (s *ProductService) CheckProduct(id int) bool {
	return s.repos.CheckProductPostgres(id)
}

func (s *ProductService) ProductDetail(id int) (productservice.ProductDetailSerailizer, error) {
	cacheData, err := s.redisDB.Get(context.Background(), fmt.Sprintf("product%d", id)).Result()

	if err == nil {
		var product productservice.ProductDetailSerailizer
		if err := json.Unmarshal([]byte(cacheData), &product); err == nil {

			return product, nil
		}
	}

	product, err := s.repos.ProductDetailPostgres(id)
	if err != nil {
		return productservice.ProductDetailSerailizer{}, err
	}
	imgs, err := s.minioClient.PresignedListObject("product", fmt.Sprintf("product%d", id), true)
	if err != nil {
		return productservice.ProductDetailSerailizer{}, err
	}
	product.Images = append(product.Images, imgs...)

	newCache, _ := json.Marshal(product)
	if err := s.redisDB.Set(context.Background(), fmt.Sprintf("product%d", id), newCache, 2*time.Hour).Err(); err != nil {
		fmt.Printf("failed to set data, error: %s", err.Error())
	}
	return product, nil
}
