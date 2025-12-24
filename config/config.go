package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server              ServerConfig              `mapstructure:"server" json:"server"`
	UserService         UserServiceConfig         `mapstructure:"user_service" json:"user_service"`
	ProductService      PostgresServiceConfig     `mapstructure:"product_service" json:"product_service"`
	OrderService        PostgresServiceConfig     `mapstructure:"order_service" json:"order_service"`
	PaymentService      PostgresServiceConfig     `mapstructure:"payment_service" json:"payment_service"`
	InventoryService    InventoryServiceConfig    `mapstructure:"inventory_service" json:"inventory_service"`
	NotificationService NotificationServiceConfig `mapstructure:"notification_service" json:"notification_service"`
}

type ServerConfig struct {
	Port     string `mapstructure:"port" json:"port"`
	GrpcPort string `mapstructure:"grpc_port" json:"grpc_port"`
	Env      string `mapstructure:"env" json:"env"`
}

// Common Postgres Config
type PostgresConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     string `mapstructure:"port" json:"port"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
	DBName   string `mapstructure:"dbname" json:"dbname"`
	SSLMode  string `mapstructure:"sslmode" json:"sslmode"`
}

// Common Redis Config
type RedisConfig struct {
	Addr     string `mapstructure:"addr" json:"addr"`
	Password string `mapstructure:"password" json:"password"`
	DB       int    `mapstructure:"db" json:"db"`
}

// Common Mongo Config
type MongoConfig struct {
	URI      string `mapstructure:"uri" json:"uri"`
	Database string `mapstructure:"database" json:"database"`
}

// Service Specific Config Wrappers
type PostgresServiceConfig struct {
	DB PostgresConfig `mapstructure:"db" json:"db"`
}

type InventoryServiceConfig struct {
	DB    PostgresConfig `mapstructure:"db" json:"db"`
	Redis RedisConfig    `mapstructure:"redis" json:"redis"`
}

type UserServiceConfig struct {
	DB    PostgresConfig `mapstructure:"db" json:"db"`
	Redis RedisConfig    `mapstructure:"redis" json:"redis"`
}

type NotificationServiceConfig struct {
	Mongo MongoConfig `mapstructure:"mongo" json:"mongo"`
}

func LoadConfig() (*Config, error) {
	// 1. Load config.yaml (Base config)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Failed to read config.yaml: %v", err)
	}

	// 2. Load .env (Overrides)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	// Enable Environment Variable Overrides
	// e.g. USER_SERVICE.DB.HOST will map to USER_SERVICE_DB_HOST
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.MergeInConfig(); err != nil {
		log.Printf("Info: No .env file found or failed to merge, relying on defaults/env: %v", err)
	}

	// Manually bind specific env vars if needed
	// But config.yaml provides the structure, so Unmarshal should mostly work.

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
