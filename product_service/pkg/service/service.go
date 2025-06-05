package service

import (
	productservice "product_service"
	"product_service/pkg/repository"

	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
)

const (
	RedisCategory = "cache_category"
)

type Admin interface {
	CreateCategory(data productservice.CategorySerializer) (int, error)
	CreateProduct(data productservice.AdminCreateProductSerializer) (int, error)
	AddImage(product int, data map[string]productservice.FileUploadSerializer) (map[string]string, error)
}

type Product interface {
	CatregoList()([]productservice.CategorySerializer, error)
}

type JWTManager interface{   
	Parse(accessToken string) (productservice.AuthMiddlewareSerializer, error)
}

type Service struct {
	Admin
	Product
}

type Deps struct{
	Repos *repository.Repository
	JWTManager JWTManager
	MinIO *minio.Client	
	Redis *redis.Client
}

func NewService(deps Deps) *Service{
	return &Service{
		Admin: NewAdminService(deps.Repos.Admin, deps.MinIO),
		Product: MewProductService(deps.Repos.Product, deps.MinIO, deps.Redis),
	}
}