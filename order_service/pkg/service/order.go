package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	orderservice "order_service"

	"order_service/pkg/repository"
	"strings"
	"sync"
	"time"

	"github.com/IBM/sarama"
	grpc_order_service "github.com/alexkhub/OnlineStoreProto/gen/go/order_service"
	"github.com/redis/go-redis/v9"
)

type OrderService struct {
	repos       repository.Order
	redisDB     *redis.Client
	producer    sarama.SyncProducer
	gRPCProduct grpc_order_service.ProductClient
	gRPCAuth    grpc_order_service.AuthClient
}

func NewOrderService(repos repository.Order, redisDB *redis.Client, producer sarama.SyncProducer, gRPCProduct grpc_order_service.ProductClient, gRPCAuth grpc_order_service.AuthClient) *OrderService {
	return &OrderService{
		repos:       repos,
		redisDB:     redisDB,
		producer:    producer,
		gRPCProduct: gRPCProduct,
		gRPCAuth:    gRPCAuth,
	}
}

func (s *OrderService) PaymentMethodeList() ([]orderservice.PaymentMethodeSerializer, error) {
	return s.repos.PaymentMethodeListPostgres()
}

func (s *OrderService) CreateOrder(order_data orderservice.CreateOrderSerializer) (int, error) {
	id, err := s.repos.CreateOrderPostgres(order_data)
	if err != nil {
		return 0, err
	}
	go func(producer sarama.SyncProducer, kafka_data orderservice.CreateOrderKafkaMessage) {
		err := SendCreateOrderKafkaMessage(producer, kafka_data)
		if err != nil {
			log.Printf("kafka error: %v", err)
			return
		}
		log.Printf("kafka ok: create order %d", kafka_data.Id)

	}(s.producer, orderservice.CreateOrderKafkaMessage{Id: id, User: order_data.User})
	return id, nil
}

func (s *OrderService) OrderDetail(ctx context.Context, orderId int64) (orderservice.EmployeeOrderDataSerializer, error) {
	var wg sync.WaitGroup

	errCh := make(chan error, 2)
	resultCh := make(chan interface{}, 2)
	serviceCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	orderDataDB, err := s.repos.GetOrderPostgres(serviceCtx, orderId)
	if err != nil {
		return orderservice.EmployeeOrderDataSerializer{}, err
	}
	// if orderDataDB.Status == "Delivered" || orderDataDB.Status == "Canceled" {
	// 	return orderservice.EmployeeOrderDataSerializer{}, fmt.Errorf("error: order status = %s", orderDataDB.Status)
	// }

	wg.Add(2)

	go func(ctx context.Context, ctxCancel context.CancelFunc, wg *sync.WaitGroup, orderId int64) {
		defer wg.Done()

		select {
		case <-ctx.Done():
			errCh <- ctx.Err()
			return
		default:
		}

		orderPoint, err := s.repos.GetOrderPointPostgres(ctx, orderId)
		if err != nil {
			errCh <- err
			ctxCancel()
			return
		}
		productListIds := make([]int64, len(orderPoint))

		for i := range orderPoint {
			productListIds[i] = orderPoint[i].ProductId
		}
		productListNames, err := s.gRPCProduct.GetProductName(ctx, &grpc_order_service.ProductIdRequest{Id: productListIds})
		if err != nil {
			errCh <- err
			ctxCancel()
			return
		}
		for i := range orderPoint {
			orderPoint[i].ProductName = productListNames.Data[i].GetName()
		}
		resultCh <- orderPoint
	}(serviceCtx, cancel, &wg, orderId)

	go func(ctx context.Context, ctxCancel context.CancelFunc, wg *sync.WaitGroup, userId int64) {
		defer wg.Done()

		select {
		case <-ctx.Done():
			errCh <- ctx.Err()
			ctxCancel()
			return
		default:
		}
		userData, err := s.gRPCAuth.GetUserData(context.Background(), &grpc_order_service.UserResponse{Id: []int64{userId}})
		if err != nil {
			errCh <- err
			ctxCancel()
			return
		}
		if len(userData.Data) == 0 {
			errCh <- fmt.Errorf("user list is empty")
		}
		resultCh <- orderservice.OrderUserSerializer{Id: userData.Data[0].GetId(), FullName: userData.Data[0].GetFullName(), Email: userData.Data[0].Email}

	}(serviceCtx, cancel, &wg, orderDataDB.UserId)

	wg.Wait()
	close(errCh)
	close(resultCh)

	var errs []string
	for opErr := range errCh {
		errs = append(errs, opErr.Error())
	}

	if len(errs) != 0 {
		return orderservice.EmployeeOrderDataSerializer{}, fmt.Errorf("errors: %s", strings.Join(errs, ", "))
	}

	result := orderservice.EmployeeOrderDataSerializer{
		Id:             orderDataDB.Id,
		FullPrice:      orderDataDB.FullPrice,
		PaymentMethod:  orderDataDB.DeliveryMethod,
		Status:         orderDataDB.Status,
		DeliveryMethod: orderDataDB.DeliveryMethod,
		CreateAt:       orderDataDB.CreateAt,
		DeliveryDate:   orderDataDB.DeliveryDate,
	}

	for data := range resultCh {
		switch v := data.(type) {
		case []orderservice.OrderPointSerializer:
			result.OrderPoints = v
		case orderservice.OrderUserSerializer:
			result.User = v
		default:
			return orderservice.EmployeeOrderDataSerializer{}, fmt.Errorf("error: bad types")
		}
	}
	return result, nil
}


func (s *OrderService) CheckOrderPermission(ctx context.Context, orderData orderservice.OrderPermission) orderservice.MyError{
	select {
		case <-ctx.Done():
			return orderservice.MyError{Error: ctx.Err(), Code: http.StatusOK, }
		default:
	}
	err := s.repos.CheckOrderPermissionPostgres(ctx, orderData)
	if err != nil{
		return orderservice.MyError{
			Error: fmt.Errorf("you don't have access"),
			Code: http.StatusForbidden,
		}
	}
	return orderservice.MyError{}
}

func (s *OrderService) UserOrders(userId int) ([]orderservice.UserOrderListSerializer, error){
	return s.repos.UserOrdersPostgres(userId)
}