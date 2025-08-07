package repository

import (
	orderservice "order_service"

	grpc_order_service "github.com/alexkhub/OnlineStoreProto/gen/go/order_service"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

const (
	CartTable = "cart"
	OrderTable = "user_order"
	PaymentMethodeTable = "payment_method"
	OrderPointTable = "order_point"
	OrderOrderPointTable = "order_order_point"
)
type Cart interface{
	CartListPostgres(user_id int)([]orderservice.CartPostgresSerializer, error)
	CreateCartPostgres(user_id int, product_id int64)(orderservice.CartSerializer, error)
	UserCartPermissionPostgres(user_id, cart_id int ) bool
	UpdateCartPostgres(cart_id, amount int) error
	CleanCartPostgres(user_id int) error
	RemoveCartPointPostgres(cart_id int) error
}

type Order interface{
	PaymentMethodeListPostgres()([]orderservice.PaymentMethodeSerializer, error)
	CreateOrderPostgres(order_data orderservice.CreateOrderSerializer) (int, error)
}

type Admin interface {
	RemoveCartPointPostgres(product_id int) error
}

type Repository struct{
	Cart
	Order
	Admin
}

type ReposDeps struct {
	DB    *sqlx.DB
	Redis *redis.Client
	GRPCProduct grpc_order_service.ProductClient
}

func NewRepository(deps ReposDeps) *Repository {
	return &Repository{
		Cart: NewCartPotgres(deps.DB),
		Order: NewOrderPostgres(deps.DB, deps.Redis, deps.GRPCProduct),
		Admin: NewAdminPostgres(deps.DB, deps.Redis),

	}
}
 