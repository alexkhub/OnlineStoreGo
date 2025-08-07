package service

import (
	"encoding/json"
	"fmt"
	"os"

	"log"
	"net/smtp"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
	"gopkg.in/gomail.v2"
)

type SendConfirmKafkaMessageData struct {
	Id int `json:"id"`
}

func SendEmail(from, password, user_email, subject, body string) error {

	msg := "From: " + from + "\n" +
		"To: " + user_email + "\n" +
		"Subject:" + subject + "\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, password, "smtp.gmail.com"),
		from, []string{user_email}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return err
	}
	log.Println("Successfully sended to " + user_email)
	return nil
}

func SendEmailV2(from, password, user_email, subject, body, image string)  error{
	msg := gomail.NewMessage()													
    msg.SetHeader("From", from)
    msg.SetHeader("To", user_email)
    msg.SetHeader("Subject", subject)

    msg.SetBody("text/plain", body)
	if image != ""{
		msg.Attach(image)
	}
    dialer := gomail.NewDialer("smtp.gmail.com", 587, from, password)
 
    if err := dialer.DialAndSend(msg); err != nil {
		return err
	}
	return nil 
    
}


func QRGeneration(url, qr_path string) error {
	if url == "" || qr_path == ""{
		return fmt.Errorf("url or qr_path is empty")
	}
	qrc, err := qrcode.New(url)
	if err != nil {
		return fmt.Errorf("could not generate QRCode: %v", err)
	}
	
	// fmt.Sprintf("./qr_codes/%s.jpeg", qr_name)
	w, err := standard.New(qr_path)
	if err != nil {
		
		return fmt.Errorf("standard.New failed: %v", err)
	}
	
	if err = qrc.Save(w); err != nil {
		return fmt.Errorf("could not save image: %v", err)
	}
	return nil
}


func EnsureDir(path string, mode os.FileMode) error{
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, mode)
		if err != nil {
			return err
		}	
	} else if err != nil {
		return err
	}
	return nil
}

func SendConfirmKafkaMessage(producer sarama.SyncProducer, id int) error {
	var data SendConfirmKafkaMessageData
	data.Id = id
	requestID := uuid.New().String()

	userJson, err := json.Marshal(data)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: ConfirmTopic,
		Key:   sarama.StringEncoder(requestID),
		Value: sarama.StringEncoder(userJson),
	}
	_, _, err = producer.SendMessage(msg)
	return err
}


