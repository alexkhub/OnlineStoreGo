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
	ProductDelete(product_id int) error
	AdminProductDetail(id int) (productservice.AdminProductDetailSerailizer, error)
	RemoveImage(product_id int, name string) error
	UpdateProduct(product_id int, product_data productservice.AdminUpdateProductSerializer) error
}

type Product interface {
	CategoryList() ([]productservice.CategorySerializer, error)
	ProductList() ([]productservice.ProductListSerailizer, error)
	CheckProduct(id int) bool
	ProductDetail(id int) (productservice.ProductDetailSerailizer, error)
}

type JWTManager interface {
	Parse(accessToken string) (productservice.AuthMiddlewareSerializer, error)
}

type Service struct {
	Admin
	Product
}

type Deps struct {
	Repos      *repository.Repository
	JWTManager JWTManager
	MinIO      *minio.Client
	Redis      *redis.Client
}

func NewService(deps Deps) *Service {
	return &Service{
		Admin:   NewAdminService(deps.Repos.Admin, deps.MinIO, deps.Repos.MinIO, deps.Redis),
		Product: MewProductService(deps.Repos.Product, deps.Redis, deps.Repos.MinIO),
	}
}
