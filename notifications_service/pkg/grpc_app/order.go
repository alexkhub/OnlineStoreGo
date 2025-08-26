package grpcapp

import (
	"context"
	"fmt"
	"log"
	notificationsservice "notifications_service"
	"notifications_service/pkg/service"

	grpc_order_service "github.com/alexkhub/OnlineStoreProto/gen/go/order_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderGRPCService struct {
	grpc_order_service.UnimplementedNotificationServer
	service service.GRPC
}

func NewOrderGRPCService(gRPC *grpc.Server, service service.GRPC) {
	grpc_order_service.RegisterNotificationServer(gRPC, &OrderGRPCService{service: service})
}

func (g *OrderGRPCService) CheckCode(ctx context.Context, request *grpc_order_service.CodeRequest) (*grpc_order_service.CodeResponse, error) {
	if request.Code < 100000 || request.Code > 1000000 {
		return nil, status.Errorf(codes.OutOfRange, "code out of range")
	}
	allowed, err := g.service.CheckCode(ctx, notificationsservice.CheckCodeSeralizer{Code: request.Code, OrderId: request.OrderId})

	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("gRPC error: %s", err.Error()))
	}

	if !allowed {
		err = g.service.GenerateNewCode(ctx, request.OrderId)
		if err != nil {
			log.Println(err)
			return nil, status.Errorf(codes.Internal, fmt.Sprintf("gRPC error: %s", err.Error()))
		}
		return nil, status.Error(codes.ResourceExhausted, "gRPC error: confirmation time has expired")

	}
	return &grpc_order_service.CodeResponse{}, nil

}
