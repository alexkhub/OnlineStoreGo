package main

import (
	"encoding/json"
	"log"
	"notifications_service"
	"os"
	"os/signal"
	"syscall"

	"notifications_service/pkg/handlers"
	"notifications_service/pkg/repository"
	"notifications_service/pkg/service"

	"github.com/IBM/sarama"
	

	_ "github.com/lib/pq"
)

func main() {
	db, err := repository.NewDBConnect()

	if err != nil {
		log.Fatalln("db err")
	}

	repos := repository.NewRepository(repository.ReposDebs{DB: db})

	sarama_config := sarama.NewConfig()

	sarama_config.Consumer.Return.Errors = true

	producer, err := sarama.NewSyncProducer([]string{"kafka:9092"}, nil)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	consumer, err := sarama.NewConsumer([]string{"kafka:9092"}, nil)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer consumer.Close()

	services := service.NewService(service.Deps{
		Repos:    repos,
		Consumer: consumer,
		Producer: producer,
	})
	
	partConsumer, err := consumer.ConsumePartition(service.AuthTopic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to consume partition: %v", err)
	}
	defer partConsumer.Close()

	my_handlers := handlers.NewHandler(services)

	go func() {
		if err := my_handlers.InitRouter().Run(":8082"); err != nil {
			log.Fatalf("server dont start")
		}
	}()
	go func() {
		for {
			select {
			// (обработка входящего сообщения и отправка ответа в Kafka)
			case msg, ok := <-partConsumer.Messages():
				if !ok {
					log.Println("Channel closed, exiting")
					return
				}
				// Десериализация входящего сообщения из JSON
				var receivedMessage notificationsservice.AuthRegistrationResponseSerializer
				err := json.Unmarshal(msg.Value, &receivedMessage)

				if err != nil {
					log.Printf("Error unmarshaling JSON: %s", err)
					continue
				}

				message, err:= services.CreateVerifyLink(receivedMessage.Id)

				if err!= nil {
					log.Printf("Create Verify link: %s", err)
				}
				
				service.SendEmail(receivedMessage.Email, "Verify Email", message)
				log.Printf("Received message: %+v\n", receivedMessage)

			}
		}
	}()

	log.Print("NotificationService Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Print("NotificationService Shutting Down")

	if err := db.Close(); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}
