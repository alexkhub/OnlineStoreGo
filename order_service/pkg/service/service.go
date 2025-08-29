package service

import (
	"context"
	orderservice "order_service"
	"order_service/pkg/repository"

	"github.com/IBM/sarama"
	grpc_order_service "github.com/alexkhub/OnlineStoreProto/gen/go/order_service"
	"github.com/redis/go-redis/v9"
)

const (
	DeleteProductTopik = "delete_prod_topik"
	CreateOrderTopik   = "create_order_topik"
)

type Cart interface {
	CartList(user_id int) ([]orderservice.CartSerializer, error)
	CreateCart(user_id int, product_id int64) (orderservice.CartSerializer, error)
	UserCartPermission(user_id, cart_id int) bool
	UpdateCart(cart_id, amount int) error
	CleanCart(user_id int) error
	RemoveCartPoint(cart_id int) error
}

type Order interface {
	PaymentMethodeList() ([]orderservice.PaymentMethodeSerializer, error)
	CreateOrder(order_data orderservice.CreateOrderSerializer) (int, error)
	OrderDetail(ctx context.Context, orderId int64) (orderservice.EmployeeOrderDataSerializer, error)
	CheckOrderPermission(ctx context.Context, orderData orderservice.OrderPermission) orderservice.MyError
	UserOrders(userId int) ([]orderservice.UserOrderListSerializer, error)
	OrdersStatistic(userId int)(orderservice.UserOrderStatisticSerializer, error)
}


type Admin interface {
	RemoveCartPoint(product_id int) error
	OrderList(filter orderservice.OrderFilter)([]orderservice.AdminOrderListSerializer, error)
	OrdersStatistic(filter orderservice.OrderFilter)([]orderservice.AdminOrderStatisticSerializer, error)
}

type Employee interface {
	ConfirmOrderStep1(ctx context.Context, confirmData orderservice.ConfirmOrderStep1Serializer) error
	ConfirmOrderStep2(confirmData orderservice.UpdateListOrderPointSerializer) error
	ConfirmOrderStep3(ctx context.Context, confirmData orderservice.ConfirmOrderStep3Serializer) error
}

type Service struct {
	Cart
	Order
	Admin
	Employee
}

type JWTManager interface {
	Parse(accessToken string) (orderservice.AuthMiddlewareSerializer, error)
}

type Deps struct {
	Repos            *repository.Repository
	Redis            *redis.Client
	GRPCProduct      grpc_order_service.ProductClient
	GRPCNotification grpc_order_service.NotificationClient
	GRPCAuth         grpc_order_service.AuthClient
	Produces         sarama.SyncProducer
}

func NewService(deps Deps) *Service {
	return &Service{
		Cart:     NewCartService(deps.Repos.Cart, deps.GRPCProduct),
		Order:    NewOrderService(deps.Repos.Order, deps.Redis, deps.Produces, deps.GRPCProduct, deps.GRPCAuth),
		Admin:    NewAdminService(deps.Repos.Admin, deps.GRPCAuth),
		Employee: NewEmployeeService(deps.Repos.Employee, deps.GRPCNotification, deps.GRPCProduct, deps.GRPCAuth),
	}
}
