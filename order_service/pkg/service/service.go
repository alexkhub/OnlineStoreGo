package service

import (
	orderservice "order_service"
	"order_service/pkg/repository"

	grpc_order_service "github.com/alexkhub/OnlineStoreProto/gen/go/order_service"
	"github.com/redis/go-redis/v9"
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

}

type  Service struct{
	Cart
	Order
}



type JWTManager interface {
	Parse(accessToken string) (orderservice.AuthMiddlewareSerializer, error)
}

type Deps struct {
	Repos       *repository.Repository
	Redis       *redis.Client
	GRPCProduct grpc_order_service.ProductClient
	
}

func NewService(deps Deps) *Service{
	return &Service{
		Cart: NewCartService(deps.Repos.Cart, deps.GRPCProduct),
		Order: NewOrderService(deps.Repos.Order, deps.Redis),
	}
}

