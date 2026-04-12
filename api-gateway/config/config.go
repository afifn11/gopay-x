package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	JWT      JWTConfig
	Services ServicesConfig
}

type AppConfig struct {
	Name string
	Port string
	Env  string
}

type JWTConfig struct {
	AccessSecret string
}

type ServicesConfig struct {
	AuthService        string
	UserService        string
	WalletService      string
	PaymentService     string
	TransactionService string
	NotificationService string
	FraudService       string
	AuditService       string
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
			Env:  viper.GetString("APP_ENV"),
		},
		JWT: JWTConfig{
			AccessSecret: viper.GetString("JWT_ACCESS_SECRET"),
		},
		Services: ServicesConfig{
			AuthService:         viper.GetString("AUTH_SERVICE_URL"),
			UserService:         viper.GetString("USER_SERVICE_URL"),
			WalletService:       viper.GetString("WALLET_SERVICE_URL"),
			PaymentService:      viper.GetString("PAYMENT_SERVICE_URL"),
			TransactionService:  viper.GetString("TRANSACTION_SERVICE_URL"),
			NotificationService: viper.GetString("NOTIFICATION_SERVICE_URL"),
			FraudService:        viper.GetString("FRAUD_SERVICE_URL"),
			AuditService:        viper.GetString("AUDIT_SERVICE_URL"),
		},
	}
}