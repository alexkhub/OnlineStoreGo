package grpcapp

import (
	"auth_service/pkg/service"
	"context"

	"github.com/alexkhub/OnlineStoreProto/gen/go/comment"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CommentGRPCServer struct {
	grpc_product_service.UnimplementedCommentServer
	service service.GRPC
}

func NewCommentGRPCServer(gRPC *grpc.Server, service service.GRPC) {
	grpc_product_service.RegisterCommentServer(gRPC, &CommentGRPCServer{service: service})
}

func (g *CommentGRPCServer) GetUserData(ctx context.Context, request *grpc_product_service.CommentIdRequest) (*grpc_product_service.UserDataResponse, error) {
	if len(request.Id) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Id list is empty")
	}

	return g.service.GetUserData(request.GetId())
}
