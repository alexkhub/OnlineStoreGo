package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"product_service/configs"
	"product_service/pkg/handlers"
	"product_service/pkg/repository"
	"product_service/pkg/service"
	"strconv"
	"syscall"

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

	db, err := repository.NewDBConnect(dbConfig["host"], db_conf_port, dbConfig["user"], dbConfig["password"], dbConfig["dbname"], dbConfig["sslmode"])

	if err != nil {
		log.Fatalln("db err", err.Error())
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

	repos := repository.NewRepository(repository.ReposDeps{DB: db, Redis: redisdb, MinIO: minios3})
	jwt_manager := service.NewManager(viper.GetString("singing_key"))
	services := service.NewService(service.Deps{Repos: repos, JWTManager: jwt_manager, MinIO: minios3, Redis: redisdb})
	my_handlers := handlers.NewHandler(services, jwt_manager)

	go func() {
		if err := my_handlers.InitRouter().Run(":8083"); err != nil {
			log.Fatalf("server didn't start")
		}
	}()

	log.Println("ProductService Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("ProductService Shutting Down")

	if err := db.Close(); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}
