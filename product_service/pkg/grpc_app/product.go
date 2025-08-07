package grpcapp

import (
	"context"
	"product_service/pkg/service"

	grpc_order_service "github.com/alexkhub/OnlineStoreProto/gen/go/order_service"
	status "google.golang.org/grpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type ProductGRPCService struct{
	grpc_order_service.UnimplementedProductServer
	service service.GRPC
}

func NewProductGRPCServer(gRPC *grpc.Server, service service.GRPC) {
	grpc_order_service.RegisterProductServer(gRPC, &ProductGRPCService{service: service})
}

func (g *ProductGRPCService)GetProductCreateCart(ctx context.Context, request *grpc_order_service.ProductIdCreateCartRequest) (*grpc_order_service.ProductDataCreateCartResponse, error) {

	if request.Id == 0{
		return nil, status.Errorf(codes.InvalidArgument, "Id is empty")
	}
	data, err := g.service.GetProductCreateCart(ctx, request.Id)

	if err != nil{
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return data, nil
	
	
}

func (g *ProductGRPCService) GetProduct(ctx context.Context, request *grpc_order_service.ProductIdRequest) (*grpc_order_service.ProductDataResponse, error) {
	if len(request.Id) == 0{
		return nil, status.Errorf(codes.InvalidArgument, "Id is empty")
	}
	data, err := g.service.GetProduct(ctx, request.Id)

	if err != nil{
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return data, nil
}

func (g *ProductGRPCService) GetProductPrice(ctx context.Context, request *grpc_order_service.ProductIdRequest) (*grpc_order_service.ProductPriceResponse, error) {

	if len(request.Id) == 0{
		return nil, status.Errorf(codes.InvalidArgument, "Id is empty")
	}
	data, err := g.service.GetProductPrice(ctx, request.Id)

	if err != nil{
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return data, nil
}

func (g *ProductGRPCService) GetProductName(ctx context.Context, request *grpc_order_service.ProductIdRequest) (*grpc_order_service.ProductNameResponse, error) {
	
	if len(request.Id) == 0{
		return nil, status.Errorf(codes.InvalidArgument, "Id is empty")
	}
	data, err := g.service.GetProductName(ctx, request.Id)

	if err != nil{
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return data, nil
}