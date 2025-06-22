package configs

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DbConfig    `mapstructure:"db"`
	KafkaConfig `mapstructure:"kafka"`
	AppHost     string `mapstructure:"app_host"`
}

type DbConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Dbname   string `mapstructure:"dbname"`
	Sslmode  string `mapstructure:"sslmode"`
}

type KafkaConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

var AppConfig Config

func LoadConfig() {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".././configs")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Load Config error %s", err.Error())
	}
	viper.Unmarshal(&AppConfig)

}
