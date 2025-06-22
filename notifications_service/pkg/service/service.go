package service

import (
	notificationsservice "notifications_service"
	"notifications_service/pkg/repository"

	"github.com/IBM/sarama"
)

const (
	ConfirmTopic = "confirm_topic"
	AuthTopic    = "auth_topic"
	BlockTopic   = "block_topik"
	Host         = "localhost"
	Port         = ":8082"
)

type Email interface {
	SendVerifyEmail(user_email string, subject string, body string)
	CreateVerifyLink(user int) (string, error)
	AccountConfirm(uuid string) error
	SendBlockEmail(data notificationsservice.UserBlockResponseSerializer)
}

type Deps struct {
	Repos    *repository.Repository
	Consumer sarama.Consumer
	Producer sarama.SyncProducer
	From     string
	Password string
}

type Service struct {
	Email
}

func NewService(deps Deps) *Service {
	new_email_service := NewEmailService(deps.Repos.Email, deps.Consumer, deps.Producer, deps.From, deps.Password)
	return &Service{
		Email: new_email_service,
	}
}
