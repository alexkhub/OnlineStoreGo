package main

import (
	// authservice "auth_service"
	// "context"

	authservice "auth_service"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"auth_service/pkg/handlers"
	"auth_service/pkg/repository"
	"auth_service/pkg/service"

	"github.com/IBM/sarama"

	_ "github.com/lib/pq"
)

const (
	SigningKey = "fdljdcsdcsv232e3cdjif"
    SigningKey2 = "fdvsgf34$%MJP&(^JGTOIOI)"
)


func main() {
	
	db, err := repository.NewDBConnect()
	
	if err != nil{
		log.Fatalln("db err")
	}

	minios3, err := repository.NewMinIOConnect()
	
	if err != nil{
		log.Fatalln("minio err", err.Error())
	}
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

	repos := repository.NewRepository(repository.ReposDebs{DB: db, MinIO: minios3})
	
	jwt_manager := service.NewManager(SigningKey, SigningKey2)

	services := service.NewService(service.Deps{
		Repos: repos,
		JWTManager: jwt_manager,
		Producer: producer,
	})
	partConsumer, err := consumer.ConsumePartition(service.ConfirmTopic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to consume partition: %v", err)
	}
	defer partConsumer.Close()

	my_handlers := handlers.NewHandler(services)

	go func() {
		if err := my_handlers.InitRouter().Run(":8081"); err != nil {
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
				var receivedMessage authservice.ConfirmUserSerializer
				err := json.Unmarshal(msg.Value, &receivedMessage)

				if err != nil {
					log.Printf("Error unmarshaling JSON: %s", err)
					continue
				}
				 err= services.ActivateUser(receivedMessage.Id)

				if err!= nil {
					log.Printf("Activate user id=%s error: %s",receivedMessage.Id, err)
				}		
				log.Printf("Received message: %+v\n", receivedMessage)
			}
		}
	}()

	log.Print("AuthService Started")
	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Print("AuthService Shutting Down")

	if err := db.Close(); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}

}
