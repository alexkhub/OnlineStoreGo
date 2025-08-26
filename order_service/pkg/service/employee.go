package service

import (
	"context"
	orderservice "order_service"
	"order_service/pkg/repository"

	grpc_order_service "github.com/alexkhub/OnlineStoreProto/gen/go/order_service"
)

type EmployeeService struct {
	repos            repository.Employee
	gRPCNotification grpc_order_service.NotificationClient
	gRPCProduct      grpc_order_service.ProductClient
	gRPCAuth         grpc_order_service.AuthClient
}

func NewEmployeeService(repos repository.Employee, gRPCNotification grpc_order_service.NotificationClient, gRPCProduct grpc_order_service.ProductClient, gRPCAuth grpc_order_service.AuthClient) *EmployeeService {
	return &EmployeeService{
		repos:            repos,
		gRPCNotification: gRPCNotification,
		gRPCProduct:      gRPCProduct,
		gRPCAuth:         gRPCAuth,
	}
}

func (s *EmployeeService) ConfirmOrderStep1(ctx context.Context, confirmData orderservice.ConfirmOrderStep1Serializer) error {
	_, err := s.gRPCNotification.CheckCode(ctx, &grpc_order_service.CodeRequest{Code: confirmData.Code, OrderId: confirmData.Order})
	return err
}

func (s *EmployeeService) ConfirmOrderStep2(confirmData orderservice.UpdateListOrderPointSerializer) error {
	if len(confirmData.Data) == 0 {
		return nil
	}
	return s.repos.UpdateOrderPointsPosrgres(confirmData)
}


func (s *EmployeeService) ConfirmOrderStep3(ctx context.Context, confirmData orderservice.ConfirmOrderStep3Serializer) error {
	return s.repos.ConfirmOrderStep3Postgres(ctx, confirmData)
}