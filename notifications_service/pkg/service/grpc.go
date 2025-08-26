package service

import (
	"context"
	"fmt"
	"log"
	"math"
	notificationsservice "notifications_service"
	"notifications_service/pkg/repository"
	"time"

	grpc_notifications_service "github.com/alexkhub/OnlineStoreProto/gen/go/notifications_service"
)

type GRPCService struct {
	repos    repository.GRPC
	gRPCAuth grpc_notifications_service.AuthClient
	from     string
	password string
}

func NewGRPCService(repos repository.GRPC, gRPCAuth grpc_notifications_service.AuthClient, from, password string) *GRPCService {
	return &GRPCService{
		repos:    repos,
		gRPCAuth: gRPCAuth,
		from:     from,
		password: password,
	}
}

func (s *GRPCService) CheckCode(ctx context.Context, orderData notificationsservice.CheckCodeSeralizer) (bool, error) {
	code_time, err := s.repos.CheckCodePostgres(ctx, orderData)
	if err != nil {
		return false, err
	}
	if math.Abs(time.Until(code_time).Minutes()) > 5.0 {
		return false, nil
	}

	return true, nil
}

func (s *GRPCService) GenerateNewCode(ctx context.Context, orderId int64) error {
	userId, err := s.repos.GetUserIdPostgres(ctx, orderId)
	if err != nil {
		return err
	}

	userData, err := s.gRPCAuth.GetUserEmail(context.Background(), &grpc_notifications_service.UserIdRequest{Id: userId})
	if err != nil {
		return err
	}

	code, err := s.repos.GenerateNewCodePostgres(ctx, orderId)
	if err != nil {
		return err
	}

	subject := "New confirm code"
	body := fmt.Sprintf("Your confirm code - %d", code)
	go func() {
		err := SendEmailV2(s.from, s.password, userData.Email, subject, body, "")
		log.Println(err)
	}()

	return nil

}
