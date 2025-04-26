package service

import (
	"auth_service"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
    salt = "fdgbgfd1232@$lv"
)


func HashPassword(password string) (string, error) {
    password = password + salt
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}
 
func CheckPasswordHash(password, hash string) bool {
    password = password + salt
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}



func SendVerifyKafkaMessage(producer sarama.SyncProducer, user authservice.AuthRegistrationResponseSerializer) error{
	requestID := uuid.New().String()
	userJson, err := json.Marshal(user)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: AuthTopic,
		Key: sarama.StringEncoder(requestID),
		Value: sarama.StringEncoder(userJson),
	}
	_, _, err = producer.SendMessage(msg)
	return err
}

func SendBlockKafkaMessage(producer sarama.SyncProducer, data authservice.UserBlockResponseSerializer) error{
	requestID := uuid.New().String()
	userJson, err := json.Marshal(data)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: BlockTopic,
		Key: sarama.StringEncoder(requestID),
		Value: sarama.StringEncoder(userJson),
	}
	_, _, err = producer.SendMessage(msg)
	return err

}