package repository

import (
	productservice "product_service"

	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
)

const (
	CategoryTable = "category"
	ProductTable = "product"
	ImageTable = "image"
	ProductImageTable = "product_image"
	RedisCategory = "cache_category"
)

type Admin interface{
	CreateCategoryPostgres(data productservice.CategorySerializer) (int, error)
	CreateProductPostgres(data productservice.AdminCreateProductSerializer) (int, error)
	CheckProductIdPostgres(product_id int) (bool)
	AddImagePostgres(product_id int, image string)( error)
}

type Product interface{
	CatregoListPostgres(cache bool)([]productservice.CategorySerializer, error)
	
}

type Repository struct{
	Admin
	Product
}

type ReposDeps struct{
	DB *sqlx.DB
	Redis *redis.Client
	MinIO *minio.Client	
}

func NewRepository(deps ReposDeps) *Repository{
	return &Repository{
		Admin: NewAdminPostgres(deps.DB, deps.Redis, deps.MinIO),
		Product: NewProductPostgres(deps.DB, deps.Redis),
	}
}