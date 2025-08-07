package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"log"
	"notifications_service"
	"os"
	"os/signal"
	"syscall"

	"notifications_service/configs"
	grpcapp "notifications_service/grpc_app"
	"notifications_service/pkg/handlers"
	"notifications_service/pkg/repository"
	"notifications_service/pkg/service"

	"github.com/IBM/sarama"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"
)

func main() {
	configs.LoadConfig()

	dbConfig := viper.GetStringMapString("db")
	db_conf_port, _ := strconv.Atoi(dbConfig["port"])

	kafkaConfig := viper.GetStringMapString("kafka")
	kafka_conf_port, _ := strconv.Atoi(kafkaConfig["port"])

	email_config := viper.GetStringMapString("email")

	db, err := repository.NewDBConnect(dbConfig["host"], db_conf_port, dbConfig["user"], dbConfig["password"], dbConfig["dbname"], dbConfig["sslmode"])

	if err != nil {
		log.Fatalln("db err")
	}

	repos := repository.NewRepository(repository.ReposDebs{DB: db})

	sarama_config := sarama.NewConfig()

	sarama_config.Consumer.Return.Errors = true

	producer, err := sarama.NewSyncProducer([]string{fmt.Sprintf("%s:%d", kafkaConfig["host"], kafka_conf_port)}, nil)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	consumer, err := sarama.NewConsumer([]string{fmt.Sprintf("%s:%d", kafkaConfig["host"], kafka_conf_port)}, nil)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer consumer.Close()

	gRPCClient, err := grpcapp.NewGRPCClient("auth_service", 9999)

	if err != nil {
		log.Fatalf("gRPC connect error: %v", err)
	}

	services := service.NewService(service.Deps{
		Repos:    repos,
		Producer: producer,
		GRPCAuth: gRPCClient,
		From:     email_config["from"],
		Password: email_config["password"],
	})

	partConsumer, err := consumer.ConsumePartition(service.AuthTopic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to consume partition: %v", err)
	}
	defer partConsumer.Close()

	partConsumerBlock, err := consumer.ConsumePartition(service.BlockTopic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to consume partition: %v", err)
	}
	defer partConsumerBlock.Close()

	partConsumerCreateOrder, err := consumer.ConsumePartition(service.CreateOrderTopik, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to consume partition: %v", err)
	}
	defer partConsumerCreateOrder.Close()

	my_handlers := handlers.NewHandler(services)

	go func() {
		if err := my_handlers.InitRouter().Run(":8082"); err != nil {
			log.Fatalf("server didn't start")
		}
	}()
	go func() {
		for{
			select {
			case msg, ok := <-partConsumer.Messages():
				if !ok {
					log.Println("Channel closed, exiting")
					return
				}
				var receivedMessage notificationsservice.AuthRegistrationResponseSerializer
				err := json.Unmarshal(msg.Value, &receivedMessage)

				if err != nil {
					log.Printf("Error unmarshaling JSON: %s", err)
					continue
				}

				message, err := services.CreateVerifyLink(receivedMessage.Id)

				if err != nil {
					log.Printf("error create verify link: %v", err)
					continue
				}

				err = services.SendVerifyEmail(receivedMessage.Email, "Verify Email", message)

				if err != nil {
					log.Printf("error send verify: %v", err)
				}
				log.Println("send verify")



			case msg, ok := <-partConsumerBlock.Messages():
				if !ok {
					log.Println("Channel closed, exiting")
					return
				}
				var receivedMessage notificationsservice.UserBlockResponseSerializer
				err := json.Unmarshal(msg.Value, &receivedMessage)

				if err != nil {
					log.Printf("Error unmarshaling JSON: %s", err)
					continue
				}
				err = services.SendBlockEmail(receivedMessage)
				if err != nil{
					log.Printf("error sevd block: %v", err)
				}
				log.Println("Block email send")

			case msg, ok := <- partConsumerCreateOrder.Messages():
				log.Println(111111)
				if !ok {
					log.Println("Channel closed, exiting")
					return
				}
				var receivedMessage notificationsservice.CreateOrderKafkaMessage
				err := json.Unmarshal(msg.Value, &receivedMessage)

				if err != nil {
					log.Printf("Error unmarshaling JSON: %s", err)
					continue
				}
				err = services.Order.SendQRForClient(receivedMessage)
				if err != nil {
					log.Printf("Error unmarshaling JSON: %s", err)
					continue
				}
				log.Println("Order QR send")

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
