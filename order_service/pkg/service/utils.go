package service

import (
	"encoding/json"
	orderservice "order_service"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)



func SendCreateOrderKafkaMessage(producer sarama.SyncProducer, data orderservice.CreateOrderKafkaMessage) error {
	requestID := uuid.New().String()
	orderJson, err := json.Marshal(data)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: CreateOrderTopik,
		Key:   sarama.StringEncoder(requestID),
		Value: sarama.StringEncoder(orderJson),
	}
	_, _, err = producer.SendMessage(msg)
	return err

}
