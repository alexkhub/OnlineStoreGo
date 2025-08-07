package service

import (
	"context"

	"product_service/pkg/repository"

	grpc_order_service "github.com/alexkhub/OnlineStoreProto/gen/go/order_service"
)

type GRPCService struct{
	repos       repository.GRPC
	minioClient repository.MinIO
}

func NewGRPCService (repos repository.GRPC, minioClient repository.MinIO) *GRPCService{
	return &GRPCService{
		repos: repos,
		minioClient: minioClient,
	}
}
func(s *GRPCService) GetProductCreateCart(ctx context.Context, productId int64)(*grpc_order_service.ProductDataCreateCartResponse, error){
	return s.repos.GetProductCreateCartPostgres(ctx, productId)
	
}

func(s *GRPCService) GetProduct(ctx context.Context, productIds []int64)(*grpc_order_service.ProductDataResponse, error){
	products, err := s.repos.GetProductPostgres(ctx, productIds)
	if err != nil{
		return nil, err
	}
	prod_data  := make([]*grpc_order_service.ProductData, 0, len(products))

	for _, value := range products{
		prod_data = append(prod_data, &grpc_order_service.ProductData{
			Id: value.Id,
			Price: value.Price,
			Name: value.Name,
		})
	}

	return &grpc_order_service.ProductDataResponse{Data: prod_data}, nil
}

func (s *GRPCService) GetProductPrice(ctx context.Context, productIds []int64) (*grpc_order_service.ProductPriceResponse, error) {
	products, err := s.repos.GetProductPricePostgres(ctx, productIds)
	if err != nil{
		return nil, err
	}
	prod_data  := make([]*grpc_order_service.ProductPrice, 0, len(products))

	for _, value := range products{
		prod_data = append(prod_data, &grpc_order_service.ProductPrice{
			Id: value.Id,
			Price: value.Price,
		
		})
	}
	return &grpc_order_service.ProductPriceResponse{Data: prod_data}, nil
}


func (s *GRPCService) GetProductName(ctx context.Context, productIds []int64) (*grpc_order_service.ProductNameResponse, error) {
	products, err := s.repos.GetProductNamePostgres(ctx, productIds)
	if err != nil{
		return nil, err
	}
	prod_data  := make([]*grpc_order_service.ProductName, 0, len(products))

	for _, value := range products{
		prod_data = append(prod_data, &grpc_order_service.ProductName{
			Id: value.Id,
			Name: value.Name,
		
		})
	}
	return &grpc_order_service.ProductNameResponse{Data: prod_data}, nil
}