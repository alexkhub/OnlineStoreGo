package service

import (
	"context"
	notificationsservice "notifications_service"
	"notifications_service/pkg/repository"

	"github.com/IBM/sarama"
	grpc_notifications_service "github.com/alexkhub/OnlineStoreProto/gen/go/notifications_service"
)

const (
	ConfirmTopic     = "confirm_topic"
	AuthTopic        = "auth_topic"
	BlockTopic       = "block_topik"
	Host             = "localhost"
	Port             = ":8082"
	CreateOrderTopik = "create_order_topik"
)

type Email interface {
	SendVerifyEmail(user_email string, subject string, body string) error
	CreateVerifyLink(user int) (string, error)
	AccountConfirm(uuid string) error
	SendBlockEmail(data notificationsservice.UserBlockResponseSerializer) error
}

type Order interface {
	SendQRForClient(orderData notificationsservice.CreateOrderKafkaMessage) error
	OrderConfirmStep1(uuid string) error
}

type GRPC interface {
	CheckCode(ctx context.Context, orderData notificationsservice.CheckCodeSeralizer) (bool, error)
	GenerateNewCode(ctx context.Context, orderId int64) error
}

type Deps struct {
	Repos    *repository.Repository
	GRPCAuth grpc_notifications_service.AuthClient
	Producer sarama.SyncProducer
	From     string
	Password string
}

type Service struct {
	Email
	Order
	GRPC
}

func NewService(deps Deps) *Service {

	return &Service{
		Email: NewEmailService(deps.Repos.Email, deps.Producer, deps.From, deps.Password),
		Order: NewOrderSerivce(deps.Repos.Order, deps.GRPCAuth, deps.Producer, deps.From, deps.Password),
		GRPC:  NewGRPCService(deps.Repos.GRPC, deps.GRPCAuth, deps.From, deps.Password),
	}
}
