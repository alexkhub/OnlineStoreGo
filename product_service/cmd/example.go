package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"product_service/configs"

	grpcapp "product_service/pkg/grpc_app"
	"product_service/pkg/handlers"
	"product_service/pkg/repository"
	"product_service/pkg/service"
	"strconv"
	"syscall"

	"github.com/IBM/sarama"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"

	"github.com/spf13/viper"
)

func main() {
	configs.LoadConfig()

	minioConfig := viper.GetStringMapString("minio")
	minio_conf_port, _ := strconv.Atoi(minioConfig["port"])
	minio_conf_useSSL, _ := strconv.ParseBool(minioConfig["use_ssl"])
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
	consumer, err := sarama.NewConsumer([]string{fmt.Sprintf("%s:%d", kafkaConfig["host"], kafka_conf_port)}, nil)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	minios3, err := repository.NewMinIOConnect(minioConfig["host"], minio_conf_port, minioConfig["access_key_id"], minioConfig["secret_access_key"], minio_conf_useSSL)

	if err != nil {
		log.Fatalln("minio err", err.Error())
	}

	redisdb, err := repository.NewRedisConnect(redisConfig["host"], redis_conf_port, redisConfig["password"])
	if err != nil {
		log.Fatalln("redis err", err.Error())
	}

	err = minios3.MakeBucket(context.Background(), "product", minio.MakeBucketOptions{Region: "eu-central-1", ObjectLocking: true})
	if err != nil {
		log.Println(err)
	}

	sarama_config := sarama.NewConfig()

	sarama_config.Consumer.Return.Errors = true

	partConsumerBlock, err := consumer.ConsumePartition(service.BlockTopic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to consume partition: %v", err)
	}
	defer partConsumerBlock.Close()

	grpcClient, err := grpcapp.NewGRPCClient(9999)

	if err != nil {
		log.Fatalf("gRPC connect error: %v", err)
	}

	repos := repository.NewRepository(repository.ReposDeps{DB: db, Redis: redisdb, MinIO: minios3})
	jwt_manager := service.NewManager(viper.GetString("singing_key"))
	services := service.NewService(service.Deps{Repos: repos, JWTManager: jwt_manager, MinIO: minios3, Redis: redisdb, GRPCComment: grpcClient})
	my_handlers := handlers.NewHandler(services, jwt_manager)

	go func() {
		if err := my_handlers.InitRouter().Run(":8083"); err != nil {
			log.Fatalf("server didn't start")
		}
	}()

	log.Println("ProductService Started")

	go func() {
		for {
			select {
			case msg, ok := <-partConsumerBlock.Messages():
				if !ok {
					log.Println("Channel closed, exiting")
					return
				}
				var user_id int
				err := json.Unmarshal(msg.Value, &user_id)

				if err != nil {
					log.Printf("Error unmarshaling JSON: %s", err)
					continue
				}
				services.RemoveUserComment(user_id)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("ProductService Shutting Down")

	if err := db.Close(); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}
