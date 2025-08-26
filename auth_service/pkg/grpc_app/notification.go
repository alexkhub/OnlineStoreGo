package grpcapp

import (
	"auth_service/pkg/service"
	"context"

	grpc_notifications_service "github.com/alexkhub/OnlineStoreProto/gen/go/notifications_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NotificationGRPCService struct {
	grpc_notifications_service.UnimplementedAuthServer
	service service.GRPC
}

func NewNotificationGRPCService(gRPC *grpc.Server, service service.GRPC) {
	grpc_notifications_service.RegisterAuthServer(gRPC, &NotificationGRPCService{service: service})
}

func (g *NotificationGRPCService) GetUserEmail(ctx context.Context, request *grpc_notifications_service.UserIdRequest) (*grpc_notifications_service.UserEmailResponse, error) {
	if request.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "Id  is empty")
	}
	data, err := g.service.GetUserEmail(request.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return data, nil

}
