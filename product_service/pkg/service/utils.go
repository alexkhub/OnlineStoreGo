package service

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)




func SendDeleteProductKafkaMessage(producer sarama.SyncProducer, product_id int) error {
	requestID := uuid.New().String()
	productJson, err := json.Marshal(product_id)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: DeleteProductTopik,
		Key:   sarama.StringEncoder(requestID),
		Value: sarama.StringEncoder(productJson),
	}
	_, _, err = producer.SendMessage(msg)
	return err

}