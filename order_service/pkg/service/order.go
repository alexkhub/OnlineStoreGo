package service

import (
	"log"
	orderservice "order_service"
	"order_service/pkg/repository"

	"github.com/IBM/sarama"
	"github.com/redis/go-redis/v9"
)

type OrderService struct{
	repos repository.Order
	redisDB     *redis.Client
	producer  sarama.SyncProducer
}

func NewOrderService(repos repository.Order, redisDB  *redis.Client, producer  sarama.SyncProducer) *OrderService{
	return &OrderService{
		repos: repos,
		redisDB: redisDB,
		producer: producer,
	}
}

func (s *OrderService) PaymentMethodeList()([]orderservice.PaymentMethodeSerializer, error){
	return s.repos.PaymentMethodeListPostgres()
}


func (s *OrderService) CreateOrder(order_data orderservice.CreateOrderSerializer) (int, error){
	id, err := s.repos.CreateOrderPostgres(order_data)
	if err != nil{
		return 0, err
	}
	go func(producer sarama.SyncProducer, kafka_data orderservice.CreateOrderKafkaMessage){
		err := SendCreateOrderKafkaMessage(producer, kafka_data)
		if err != nil{
			log.Printf("kafka error: %v", err)
			return
		}
		log.Printf("kafka ok: create order %d", kafka_data.Id)
		
	}(s.producer, orderservice.CreateOrderKafkaMessage{Id: id, User: order_data.User})
	return id, nil
}