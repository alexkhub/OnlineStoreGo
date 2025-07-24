package main

import (
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

	_ "github.com/lib/pq"

	"github.com/spf13/viper"
)


func main() {
	configs.LoadConfig()

	
	dbConfig := viper.GetStringMapString("db")
	db_conf_port, _ := strconv.Atoi(dbConfig["port"])

	redisConfig := viper.GetStringMapString("redis")
	redis_conf_port, _ := strconv.Atoi(redisConfig["port"])
	
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
	repos := repository.NewRepository(repository.ReposDeps{DB: db, Redis: redisdb,})
	jwt_manager := service.NewManager(viper.GetString("singing_key"))
	services := service.NewService(service.Deps{Repos: repos,  Redis: redisdb, GRPCProduct: grpcClient})
	my_handlers := handlers.NewHandler(services, jwt_manager)

	go func() {
		if err := my_handlers.InitRouter().Run(fmt.Sprintf(":%s", viper.GetString("app_host"))); err != nil {
			log.Fatalf("server didn't start")
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
