package grpcapp

import (
	"context"

	grpc_order_service "github.com/alexkhub/OnlineStoreProto/gen/go/order_service"

	"auth_service/pkg/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderGRPCService struct {
	grpc_order_service.UnimplementedAuthServer
	service service.GRPC
}

func NewOrderGRPCService(gRPC *grpc.Server, service service.GRPC) {
	grpc_order_service.RegisterAuthServer(gRPC, &OrderGRPCService{service: service})
}

func (g *OrderGRPCService) GetUserData(ctx context.Context, request *grpc_order_service.UserResponse) (*grpc_order_service.UserDataResponse, error) {
	if len(request.Id) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Id list is empty")
	}

	data, err := g.service.GetOrderUserData(request.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return data, nil

}
