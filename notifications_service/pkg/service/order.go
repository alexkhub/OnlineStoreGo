package service

import (
	"context"
	"fmt"
	notificationsservice "notifications_service"
	"notifications_service/pkg/repository"

	"github.com/IBM/sarama"
	grpc_notifications_service "github.com/alexkhub/OnlineStoreProto/gen/go/notifications_service"
)

type OrderService struct {
	repos    repository.Order
	gRPCAuth grpc_notifications_service.AuthClient
	producer sarama.SyncProducer
	from     string
	password string
}

func NewOrderSerivce(repos repository.Order, gRPCAuth grpc_notifications_service.AuthClient, producer sarama.SyncProducer, from, password string) *OrderService {
	return &OrderService{
		repos:    repos,
		gRPCAuth: gRPCAuth,
		producer: producer,
		from:     from,
		password: password,
	}
}

func (s *OrderService) SendQRForClient(orderData notificationsservice.CreateOrderKafkaMessage) error {

	uuid, err := s.repos.CreateVerifyPostgres(orderData.User, orderData.Id)
	if err != nil {
		return err
	}

	err = EnsureDir("../qr_codes", 0777)
	if err != nil {
		return err
	}

	qr_path := fmt.Sprintf("../qr_codes/order%d.jpeg", orderData.Id)
	qr_url := fmt.Sprintf("%s%s/order_qr/%s", Host, Port, uuid)

	err = QRGeneration(qr_url, qr_path)
	if err != nil {
		return err
	}

	userData, err := s.gRPCAuth.GetUserEmail(context.Background(), &grpc_notifications_service.UserIdRequest{Id: int64(orderData.User)})
	if err != nil {
		return err
	}

	subject := "Order QR"
	body := "This QR must be provided upon receipt of the order"

	err = SendEmailV2(s.from, s.password, userData.Email, subject, body, qr_path)
	if err != nil {
		return err
	}

	return nil
}

func (s *OrderService) OrderConfirmStep1(uuid string) error {

	data, err := s.repos.CheckUUIDPostgres(uuid)

	if err != nil {
		return err
	}

	userData, err := s.gRPCAuth.GetUserEmail(context.Background(), &grpc_notifications_service.UserIdRequest{Id: int64(data.UserId)})

	if err != nil {
		return err
	}

	code, err := s.repos.CodeGenerationPostgres(data.OrderId)
	if err != nil {
		return err
	}

	subject := "Confirm order"
	body := fmt.Sprintf("Your confirm code - %d", code)

	err = SendEmailV2(s.from, s.password, userData.Email, subject, body, "")
	if err != nil {
		return err
	}
	return err
}
