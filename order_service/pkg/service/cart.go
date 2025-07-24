package service

import (
	"context"
	orderservice "order_service"
	"order_service/pkg/repository"

	grpc_order_service "github.com/alexkhub/OnlineStoreProto/gen/go/order_service"
)


type CartService struct {
	repos repository.Cart
	gRPCProduct grpc_order_service.ProductClient
}

func NewCartService(repos repository.Cart, gRPCProduct grpc_order_service.ProductClient) *CartService{
	return  &CartService{
		repos: repos,
		gRPCProduct: gRPCProduct,
	}
}


func (s *CartService) CartList(user_id int)([]orderservice.CartSerializer, error){
	dataDB, err := s.repos.CartListPostgres(user_id)
	if err != nil{
		return nil, err
	}
	if len(dataDB) == 0{
		return nil, nil
	}

	productListId := make([]int64, 0, len(dataDB)) 

	for _, value := range dataDB{
		productListId = append(productListId, value.Product)
	}

	productData, err := s.gRPCProduct.GetProduct(context.Background(), &grpc_order_service.ProductIdRequest{Id: productListId})
	if err != nil{
		return nil, err
	}
	productDataMap := make(map[int64]orderservice.CartProductSerializer)

	for _, value := range productData.Data{
		productDataMap[value.Id] = orderservice.CartProductSerializer{
			Id: value.Id,
			Price: value.Price,
			Name: value.Name,

		}
	}
	cartList := make([]orderservice.CartSerializer, 0, len(dataDB))
	for _, value := range dataDB{
		cartList = append(cartList, orderservice.CartSerializer{
			Id: value.Id, 
			Amount: value.Amount, 
			Product: productDataMap[value.Product]})

	}
	return cartList, nil
	
}


func (s *CartService) CreateCart(user_id int, product_id int64)(orderservice.CartSerializer, error){

	prodData, err := s.gRPCProduct.GetProductCreateCart(context.Background(), &grpc_order_service.ProductIdCreateCartRequest{Id: product_id})
	if err != nil{
		return orderservice.CartSerializer{}, err
	}

	cartData, err := s.repos.CreateCartPostgres(user_id, product_id)
	if err != nil{
		return orderservice.CartSerializer{}, err
	}
	cartData.Product = orderservice.CartProductSerializer{Id : prodData.Id, Price: prodData.Price, Name: prodData.Name}

	return cartData, nil
}
