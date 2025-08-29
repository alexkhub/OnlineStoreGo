package service

import (
	"context"
	orderservice "order_service"
	"order_service/pkg/repository"

	grpc_order_service "github.com/alexkhub/OnlineStoreProto/gen/go/order_service"
)

type AdminService struct {
	repos repository.Admin
	gRPCAuth    grpc_order_service.AuthClient
}

func NewAdminService(repos repository.Admin, gRPCAuth grpc_order_service.AuthClient) *AdminService {
	return &AdminService{
		repos: repos,
		gRPCAuth: gRPCAuth,
	}
}

func (s *AdminService) RemoveCartPoint(product_id int) error {
	return s.repos.RemoveCartPointPostgres(product_id)
}


func (s *AdminService) OrderList(filter orderservice.OrderFilter)([]orderservice.AdminOrderListSerializer, error){
	var orders []orderservice.AdminOrderListSerializer
	ordersDB, err := s.repos.OrderListPostgres(filter)
	if err != nil{
		return nil, err
	}

	if len(ordersDB) == 0{
		return orders, nil
	}

	UniqueUserIdMap := make(map[int64]struct{},) 

	for _, value :=  range ordersDB{
		UniqueUserIdMap[value.UserId] = struct{}{}
	}

	UniqueUserId := make([]int64, 0, len(UniqueUserIdMap))

	for key := range UniqueUserIdMap{
		UniqueUserId = append(UniqueUserId, key)
	}
	userData, err := s.gRPCAuth.GetUserData(context.Background(), &grpc_order_service.UserResponse{Id:UniqueUserId})
	if err != nil{
		return nil, err
	}

	userDataMap := make(map[int64]*grpc_order_service.UserData )
	for _, value := range userData.Data{
		userDataMap[value.Id] = value
	}


	for _, value := range ordersDB{
		user := userDataMap[value.UserId]
		orders = append(orders, orderservice.AdminOrderListSerializer{
			Id: value.Id,
			UserId: orderservice.OrderUserSerializer{
				Id: user.Id,
				Email: user.Email,
				FullName: user.FullName,
			},
			FullPrice: value.FullPrice,
			DeliveryMethod: value.DeliveryMethod,
			PaymentStatus: value.PaymentStatus,
			CreateAt: value.CreateAt,
			Status: value.Status,
			DeliveryDate: value.DeliveryDate,
		})
	}

	return orders, nil
}


func (s *AdminService) OrdersStatistic(filter orderservice.OrderFilter)([]orderservice.AdminOrderStatisticSerializer, error){
	return s.repos.OrdersStatisticPostgres(filter)
}