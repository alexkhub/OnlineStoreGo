package service

import (
	orderservice "order_service"
	"order_service/pkg/repository"

	"github.com/IBM/sarama"
	grpc_order_service "github.com/alexkhub/OnlineStoreProto/gen/go/order_service"
	"github.com/redis/go-redis/v9"
)

const (
	DeleteProductTopik = "delete_prod_topik"
	CreateOrderTopik = "create_order_topik"
)


type Cart interface{
	CartList(user_id int)([]orderservice.CartSerializer, error)
	CreateCart(user_id int, product_id int64)(orderservice.CartSerializer, error)
	UserCartPermission(user_id, cart_id int ) bool
	UpdateCart(cart_id, amount int) error
	CleanCart(user_id int) error
	RemoveCartPoint(cart_id int) error
	
}

type Order interface{
	PaymentMethodeList()([]orderservice.PaymentMethodeSerializer, error)
	CreateOrder(order_data orderservice.CreateOrderSerializer) (int, error)
}

type Admin interface{
	RemoveCartPoint(product_id int) error
}

type  Service struct{
	Cart
	Order
	Admin
}



type JWTManager interface {
	Parse(accessToken string) (orderservice.AuthMiddlewareSerializer, error)
}

type Deps struct {
	Repos       *repository.Repository
	Redis       *redis.Client
	GRPCProduct grpc_order_service.ProductClient
	Produces sarama.SyncProducer	
}

func NewService(deps Deps) *Service{
	return &Service{
		Cart: NewCartService(deps.Repos.Cart, deps.GRPCProduct),
		Order: NewOrderService(deps.Repos.Order, deps.Redis, deps.Produces),
		Admin: NewAdminService(deps.Repos.Admin),
	}
}

