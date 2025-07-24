package grpcapp

import (
	"fmt"

	grpc_order_service "github.com/alexkhub/OnlineStoreProto/gen/go/order_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGRPCClient(host string, port int) (grpc_order_service.ProductClient,  error) {
	connect, err := grpc.NewClient(fmt.Sprintf("%s:%d", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	grpcClient := grpc_order_service.NewProductClient(connect)

	return grpcClient, nil

}
