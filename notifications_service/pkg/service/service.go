package service

import (
	"notifications_service/pkg/repository"
	"github.com/IBM/sarama"
)

const (
	ConfirmTopic = "confirm_topic"
	AuthTopic = "auth_topic"
	Host = "localhost"
	Port = ":8082"
)

type Email interface {
	SendVerifyEmail()()
	CreateVerifyLink(user int) (string, error)
	AccountConfirm(uuid string) (error)

}

type Deps struct {
    Repos *repository.Repository
	Consumer sarama.Consumer
	Producer sarama.SyncProducer
	
	
}

type Service struct {
	Email
	
}

func NewService(deps Deps) *Service{
	new_email_service := NewEmailService(deps.Repos.Email, deps.Consumer, deps.Producer)
	return &Service{
		Email: new_email_service,
		
	}
}