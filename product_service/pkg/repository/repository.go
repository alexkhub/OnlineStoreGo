package repository

import (
	productservice "product_service"

	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
)

const (
	CategoryTable     = "category"
	ProductTable      = "product"
	ImageTable        = "image"
	ProductImageTable = "product_image"
	CommentTable      = "comment"
	RedisCategory     = "cache_category"
)

type Admin interface {
	CreateCategoryPostgres(data productservice.CategorySerializer) (int, error)
	CreateProductPostgres(data productservice.AdminCreateProductSerializer) (int, error)
	CheckProductIdPostgres(product_id int) bool
	AddImagePostgres(product_id int, image string) error
	DeleteProductPostgres(id int) error
	AdminProductDetailPostgres(id int) (productservice.AdminProductDetailSerailizer, error)
	GetImagesPostgres(product_id int) ([]productservice.ImageSerializer, error)
	DeleteImagePostgres(name string) error
	UpdateProductPostgres(product_id int, product_data productservice.AdminUpdateProductSerializer) error

	RemoveCommentPostgres( comment_id int) error 
}

type Product interface {
	CategoryListPostgres(cache bool) ([]productservice.CategorySerializer, error)
	ProductListPostgres() ([]productservice.ProductListSerailizer, error)
	CheckProductPostgres(id int) bool
	ProductDetailPostgres(id int) (productservice.ProductDetailSerailizer, error)
}

type MinIO interface {
	GetOne(bucketName, objectID string) (string, error)
	GetMany(bucketName string, objectIDs []string) (map[string]string, error)
	PresignedListObject(bucketName, prefix string, recursive bool) ([]string, error)
	RemoveAllObjects(bucketName, prefix string, recursive bool)
	RemoveOne(bucketName, objectID string) error
}

type Comment interface {
	CreateCommentPostgres(data productservice.CreateCommentSerializer, product_id, user_id int) (int, error)
	RemoveUserCommentPostgres(user_id int) error
	CommentListPostgres(product_id int) ([]productservice.ListCommentPostgresSerializer, error)
	RemoveCommentPostgres(comment_id int, user_id int) error

}


type Repository struct {
	Admin
	Product
	MinIO
	Comment
}

type ReposDeps struct {
	DB    *sqlx.DB
	Redis *redis.Client
	MinIO *minio.Client
}

func NewRepository(deps ReposDeps) *Repository {
	return &Repository{
		Admin:   NewAdminPostgres(deps.DB, deps.Redis, deps.MinIO),
		Product: NewProductPostgres(deps.DB, deps.Redis, deps.MinIO),
		MinIO:   NewMinioClient(deps.MinIO),
		Comment: NewCommentPostgres(deps.DB),
	}
}
