package grpcapp

import (
	// "fmt"

	"fmt"

	grpc_product_service "github.com/alexkhub/OnlineStoreProto/gen/go/comment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGRPCClient(host string, port int) (grpc_product_service.CommentClient, error) {

	connect, err := grpc.NewClient(fmt.Sprintf("%s:%d", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	grpcClient := grpc_product_service.NewCommentClient(connect)

	return grpcClient, nil

}
