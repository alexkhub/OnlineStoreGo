package service

import (
	"errors"
	"fmt"
	"log"
	notificationsservice "notifications_service"
	"notifications_service/pkg/repository"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type EmailService struct {
	repos repository.Email
	consumer sarama.Consumer
	producer sarama.SyncProducer

}

func NewEmailService(repos repository.Email, consumer sarama.Consumer, producer sarama.SyncProducer) *EmailService{
	return &EmailService{
		repos:  repos,
		consumer: consumer,
		producer: producer,
	}
}

func (s *EmailService) SendVerifyEmail(user_email string, subject string, body string){
	go func (){
		err := SendEmail(user_email , subject , body)
		if err!=nil{
			log.Printf("Send verify email error %s", err)
		}
	}()
}

func (s *EmailService) CreateVerifyLink(user int) (string, error){
	uuid := uuid.New().String()
	err := s.repos.CreateVerify(uuid, user)
	if err!= nil{
		return  "", err
	}
	result := fmt.Sprintf("Follow the link to activate your account %s%s/confirm/%s", Host, Port, uuid)
	return result, nil
}

func (s *EmailService) AccountConfirm(uuid string) (error){
	time := time.Now()
	user_id, cheate_time, err := s.repos.ChechUUID(uuid)

	if err != nil{
		return err 
	}
	if int(time.Sub(cheate_time).Minutes()) > 30{
		return errors.New("Time empty")
	}
	err = SendConfirmKafkaMessage(s.producer, user_id)

	log.Println(user_id)


	return err
}

func (s *EmailService) SendBlockEmail(data notificationsservice.UserBlockResponseSerializer){
	go func(){
		var subject, body string
		
		if data.Block{
			subject = "Acccount block"
			body = fmt.Sprintf("You account blocked at %s", time.Now().Format(time.DateTime))
		}else{
			subject = "Acccount unblock"
			body = fmt.Sprintf("You account unblocked at %s", time.Now().Format(time.DateTime))
		}
		err := SendEmail(data.Email , subject , body)
		if err!=nil{
			log.Printf("Send block email error %s", err)
		}

	}()
}