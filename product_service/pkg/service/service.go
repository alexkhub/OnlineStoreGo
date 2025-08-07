package service

import (
	"context"
	productservice "product_service"
	"product_service/pkg/repository"

	"github.com/IBM/sarama"
	grpc_product_service "github.com/alexkhub/OnlineStoreProto/gen/go/product_service"
	grpc_order_service "github.com/alexkhub/OnlineStoreProto/gen/go/order_service"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
)

const (
	RedisCategory = "cache_category"
	BlockTopic    = "block_topik_v2"
	DeleteProductTopik = "delete_prod_topik"
)

type Admin interface {
	CreateCategory(data productservice.CategorySerializer) (int, error)
	CreateProduct(data productservice.AdminCreateProductSerializer) (int, error)
	AddImage(product int, data map[string]productservice.FileUploadSerializer) (map[string]string, error)
	ProductDelete(product_id int) error
	AdminProductDetail(id int) (productservice.AdminProductDetailSerailizer, error)
	RemoveImage(product_id int, name string) error
	UpdateProduct(product_id int, product_data productservice.AdminUpdateProductSerializer) error
	RemoveComment(comment_id int) error
}

type Product interface {
	CategoryList() ([]productservice.CategorySerializer, error)
	ProductList() ([]productservice.ProductListSerailizer, error)
	CheckProduct(id int) bool
	ProductDetail(id int) (productservice.ProductDetailSerailizer, error)
}

type Comment interface {
	CreateComment(data productservice.CreateCommentSerializer, product_id, user_id int) (int, error)
	RemoveUserComment(user_id int) error
	CommentList(product_id int) ([]productservice.ListCommentSerializer, error)
	RemoveComment(comment_id int, user_id int) error
}

type GRPC interface{
	GetProductCreateCart(ctx context.Context, productId int64)(*grpc_order_service.ProductDataCreateCartResponse, error)
	GetProduct(ctx context.Context, productIds []int64)(*grpc_order_service.ProductDataResponse, error)
	GetProductPrice(ctx context.Context, productIds []int64) (*grpc_order_service.ProductPriceResponse, error) 
	GetProductName(ctx context.Context, productIds []int64) (*grpc_order_service.ProductNameResponse, error)

}


type JWTManager interface {
	Parse(accessToken string) (productservice.AuthMiddlewareSerializer, error)
}

type Service struct {
	Admin
	Product
	Comment
	GRPC
}


type Deps struct {
	Repos       *repository.Repository
	MinIO       *minio.Client
	Redis       *redis.Client
	GRPCComment grpc_product_service.CommentClient
	Producer   sarama.SyncProducer
	
}

func NewService(deps Deps) *Service {
	return &Service{
		Admin:   NewAdminService(deps.Repos.Admin, deps.MinIO, deps.Repos.MinIO, deps.Redis, deps.Producer),
		Product: MewProductService(deps.Repos.Product, deps.Redis, deps.Repos.MinIO),
		Comment: NewCommentService(deps.Repos.Comment, deps.Redis, deps.GRPCComment),
		GRPC: NewGRPCService(deps.Repos.GRPC, deps.Repos.MinIO),

	}
}
