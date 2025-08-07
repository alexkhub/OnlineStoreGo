package grpcapp

import (
	"fmt"
	grpc_notifications_service "github.com/alexkhub/OnlineStoreProto/gen/go/notifications_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGRPCClient(host string, port int) (grpc_notifications_service.AuthClient,  error) {
	connect, err := grpc.NewClient(fmt.Sprintf("%s:%d", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	grpcClient := grpc_notifications_service.NewAuthClient(connect)

	return grpcClient, nil

}