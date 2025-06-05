package configs

import (
	"log"

	"github.com/spf13/viper"
)


type Config struct {

	DbConfig `mapstructure:"db"`
	MinioConfig `mapstructure:"minio"`
	AppHost string `mapstructure:"app_host"`
	SingingKey string `mapstructure:"singing_key"`
}




type DbConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	User string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Dbname string `mapstructure:"dbname"`
	Sslmode string `mapstructure:"sslmode"`

}
type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	Password string `mapstructure:"password"`
}

type MinioConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	AccessKeyID string `mapstructure:"accessKeyID"`
	SecretAccessKey string `mapstructure:"secretAccessKey"`
	UseSSL string `mapstructure:"useSSL"`

}

var AppConfig Config

func LoadConfig(){

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	// viper.AddConfigPath("./configs")
	viper.AddConfigPath(".././configs")

	
	
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Load Config error %s", err.Error())
	}
	viper.Unmarshal(&AppConfig)

}