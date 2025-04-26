package service

import (
	"encoding/json"

	"log"
	"net/smtp"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type SendConfirmKafkaMessageData struct{
    Id int `json: "id"`
}

func SendEmail(user_email string, subject string, body string) (error){
    from := "aleksandrkhubaevwork@gmail.com"
    pass := "qdfgbwcyublqpler"


    msg := "From: " + from + "\n" +
        "To: " + user_email + "\n" +
        "Subject:"+ subject + "\n\n" +
        body

    err := smtp.SendMail("smtp.gmail.com:587",
        smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
        from, []string{user_email}, []byte(msg))

    if err != nil {
        log.Printf("smtp error: %s", err)
        return err 
    }
    log.Println("Successfully sended to " + user_email)
	return nil
}



func SendConfirmKafkaMessage(producer sarama.SyncProducer, id int) (error){
    var data SendConfirmKafkaMessageData
    data.Id = id
    requestID := uuid.New().String()
    
	userJson, err := json.Marshal(data)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: ConfirmTopic,
		Key: sarama.StringEncoder(requestID),
		Value: sarama.StringEncoder(userJson),
	}
	_, _, err = producer.SendMessage(msg)
	return err
}