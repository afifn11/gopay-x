package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App   AppConfig
	Kafka KafkaConfig
}

type AppConfig struct {
	Name string
	Port string
}

type KafkaConfig struct {
	Brokers string
}

func Load() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("No .env file found, reading from environment variables")
	}

	return &Config{
		App: AppConfig{
			Name: viper.GetString("APP_NAME"),
			Port: viper.GetString("APP_PORT"),
		},
		Kafka: KafkaConfig{
			Brokers: viper.GetString("KAFKA_BROKERS"),
		},
	}
}