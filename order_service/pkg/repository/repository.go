package repository

import (
	"context"
	orderservice "order_service"

	grpc_order_service "github.com/alexkhub/OnlineStoreProto/gen/go/order_service"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

const (
	CartTable            = "cart"
	OrderTable           = "user_order"
	PaymentMethodeTable  = "payment_method"
	OrderPointTable      = "order_point"
	OrderOrderPointTable = "order_order_point"
)

type Cart interface {
	CartListPostgres(user_id int) ([]orderservice.CartPostgresSerializer, error)
	CreateCartPostgres(user_id int, product_id int64) (orderservice.CartSerializer, error)
	UserCartPermissionPostgres(user_id, cart_id int) bool
	UpdateCartPostgres(cart_id, amount int) error
	CleanCartPostgres(user_id int) error
	RemoveCartPointPostgres(cart_id int) error
}

type Order interface {
	PaymentMethodeListPostgres() ([]orderservice.PaymentMethodeSerializer, error)
	CreateOrderPostgres(order_data orderservice.CreateOrderSerializer) (int, error)
	GetOrderPostgres(ctx context.Context, orderId int64) (orderservice.EmployeePreparatoryOrderDataSerializer, error)
	GetOrderPointPostgres(ctx context.Context, orderId int64) ([]orderservice.OrderPointSerializer, error)
	CheckOrderPermissionPostgres(ctx context.Context, orderData orderservice.OrderPermission) error
	UserOrdersPostgres(userId int) ([]orderservice.UserOrderListSerializer, error)
	
}

type Admin interface {
	RemoveCartPointPostgres(product_id int) error
}

type Employee interface {
	UpdateOrderPointsPosrgres(confirmData orderservice.UpdateListOrderPointSerializer) error
	ConfirmOrderStep3Postgres(ctx context.Context, confirmData orderservice.ConfirmOrderStep3Serializer) error
}

type Repository struct {
	Cart
	Order
	Admin
	Employee
}

type ReposDeps struct {
	DB          *sqlx.DB
	Redis       *redis.Client
	GRPCProduct grpc_order_service.ProductClient
}

func NewRepository(deps ReposDeps) *Repository {
	return &Repository{
		Cart:     NewCartPotgres(deps.DB),
		Order:    NewOrderPostgres(deps.DB, deps.Redis, deps.GRPCProduct),
		Admin:    NewAdminPostgres(deps.DB, deps.Redis),
		Employee: NewEmloyeePostgres(deps.DB),
	}
}
