package main

import (
	"encoding/json"
	"fmt"
	"log"
	configs "order_service/configs"
	grpcapp "order_service/pkg/grpc_app"
	"order_service/pkg/handlers"
	"order_service/pkg/repository"
	"order_service/pkg/service"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/IBM/sarama"
	_ "github.com/lib/pq"

	"github.com/spf13/viper"
)


func main() {
	configs.LoadConfig()
	sarama_config := sarama.NewConfig()
	sarama_config.Consumer.Return.Errors = true

	
	dbConfig := viper.GetStringMapString("db")
	db_conf_port, _ := strconv.Atoi(dbConfig["port"])

	redisConfig := viper.GetStringMapString("redis")
	redis_conf_port, _ := strconv.Atoi(redisConfig["port"])

	kafkaConfig := viper.GetStringMapString("kafka")
	kafka_conf_port, _ := strconv.Atoi(kafkaConfig["port"])
	
	db, err := repository.NewDBConnect(dbConfig["host"], db_conf_port, dbConfig["user"], dbConfig["password"], dbConfig["dbname"], dbConfig["sslmode"])

	if err != nil {
		log.Fatalln("db err", err.Error())
	}


	redisdb, err := repository.NewRedisConnect(redisConfig["host"], redis_conf_port, redisConfig["password"])
	if err != nil {
		log.Fatalln("redis err", err.Error())
	}

	grpcClient, err := grpcapp.NewGRPCClient("product_service", 9999)
	if err != nil {
		log.Fatalln("gRPC product err", err.Error())
	}

	consumer, err := sarama.NewConsumer([]string{fmt.Sprintf("%s:%d", kafkaConfig["host"], kafka_conf_port)}, nil)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}

	producer, err := sarama.NewSyncProducer([]string{fmt.Sprintf("%s:%d", kafkaConfig["host"], kafka_conf_port)}, nil)
	if err != nil {
		log.Fatalf("failed to create producer: %v", err)
	}
	defer producer.Close()

	partConsumerDelProd, err := consumer.ConsumePartition(service.DeleteProductTopik, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to consume partition: %v", err)
	}
	defer partConsumerDelProd.Close()

	repos := repository.NewRepository(repository.ReposDeps{DB: db, Redis: redisdb, GRPCProduct: grpcClient})
	jwt_manager := service.NewManager(viper.GetString("singing_key"))
	services := service.NewService(service.Deps{Repos: repos,  Redis: redisdb, GRPCProduct: grpcClient, Produces: producer})
	my_handlers := handlers.NewHandler(services, jwt_manager)

	go func() {
		if err := my_handlers.InitRouter().Run(fmt.Sprintf(":%s", viper.GetString("app_host"))); err != nil {
			log.Fatalf("server didn't start")
		}
	}()
	
	go func() {
		for {
			select {
			case msg, ok := <-partConsumerDelProd.Messages():
				if !ok {
					log.Println("channel closed, exiting")
					return
				}
				var prodId int
				err := json.Unmarshal(msg.Value, &prodId)

				if err != nil {
					log.Printf("Error unmarshaling JSON: %s", err)
					continue
				}
				err = services.Admin.RemoveCartPoint(prodId)

				if err != nil {
					log.Printf("Remove product id=%d error: %s", prodId, err)
				}
				log.Printf("Received message: %+v\n", prodId)
			}
		}
	}()
	log.Println("Order Service Started")

	

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("Order Service Shutting Down")

	if err := db.Close(); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}
