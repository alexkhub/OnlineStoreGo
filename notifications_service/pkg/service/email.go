package service

import (
	"fmt"

	notificationsservice "notifications_service"
	"notifications_service/pkg/repository"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type EmailService struct {
	repos repository.Email

	producer sarama.SyncProducer
	from     string
	password string
}

func NewEmailService(repos repository.Email, producer sarama.SyncProducer, from, password string) *EmailService {
	return &EmailService{
		repos:    repos,
		producer: producer,
		from:     from,
		password: password,
	}
}

func (s *EmailService) SendVerifyEmail(user_email string, subject string, body string) error {

	return SendEmailV2(s.from, s.password, user_email, subject, body, "")

}

func (s *EmailService) CreateVerifyLink(user int) (string, error) {
	uuid := uuid.New().String()
	err := s.repos.CreateVerify(uuid, user)
	if err != nil {
		return "", err
	}
	result := fmt.Sprintf("Follow the link to activate your account %s%s/confirm/%s", Host, Port, uuid)
	return result, nil
}

func (s *EmailService) AccountConfirm(uuid string) error {
	time := time.Now()
	data, err := s.repos.ChechUUID(uuid)

	if err != nil {
		return err
	}
	if int(time.Sub(data.CreateTime).Minutes()) > 30 {
		return fmt.Errorf("the term has expired")
	}
	err = SendConfirmKafkaMessage(s.producer, data.UserId)

	return err
}

func (s *EmailService) SendBlockEmail(data notificationsservice.UserBlockResponseSerializer) error {

	var subject, body string

	if data.Block {
		subject = "Acccount block"
		body = fmt.Sprintf("You account blocked at %s", time.Now().Format(time.DateTime))
	} else {
		subject = "Acccount unblock"
		body = fmt.Sprintf("You account unblocked at %s", time.Now().Format(time.DateTime))
	}
	return SendEmail(s.from, s.password, data.Email, subject, body)

}
