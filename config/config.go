package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server              ServerConfig              `json:"server"`
	UserService         UserServiceConfig         `json:"user_service"`
	ProductService      PostgresServiceConfig     `json:"product_service"`
	OrderService        PostgresServiceConfig     `json:"order_service"`
	PaymentService      PostgresServiceConfig     `json:"payment_service"`
	InventoryService    InventoryServiceConfig    `json:"inventory_service"`
	NotificationService NotificationServiceConfig `json:"notification_service"`
}

type ServerConfig struct {
	Port     string `json:"port"`
	GrpcPort string `json:"grpc_port"`
	Env      string `json:"env"`
}

// Common Postgres Config
type PostgresConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
}

// Common Redis Config
type RedisConfig struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// Common Mongo Config
type MongoConfig struct {
	URI      string `json:"uri"`
	Database string `json:"database"`
}

// Service Specific Config Wrappers
type PostgresServiceConfig struct {
	DB PostgresConfig `json:"db"`
}

type InventoryServiceConfig struct {
	DB    PostgresConfig `json:"db"`
	Redis RedisConfig    `json:"redis"`
}

type UserServiceConfig struct {
	DB    PostgresConfig `json:"db"`
	Redis RedisConfig    `json:"redis"`
}

type NotificationServiceConfig struct {
	Mongo MongoConfig `json:"mongo"`
}

func LoadConfig() (*Config, error) {
	// 1. Load .env using godotenv
	if err := godotenv.Load(); err != nil {
		log.Printf("Info: No .env file found, relying on system vars: %v", err)
	}

	cfg := &Config{
		Server: ServerConfig{
			Port:     getEnv("SERVER_PORT", "8000"),
			GrpcPort: getEnv("SERVER_GRPC_PORT", "50051"),
			Env:      getEnv("SERVER_ENV", "development"),
		},
		UserService: UserServiceConfig{
			DB:    loadDBConfig("USER_SERVICE_DB"),
			Redis: loadRedisConfig("USER_SERVICE_REDIS"),
		},
		ProductService: PostgresServiceConfig{
			DB: loadDBConfig("PRODUCT_SERVICE_DB"),
		},
		OrderService: PostgresServiceConfig{
			DB: loadDBConfig("ORDER_SERVICE_DB"),
		},
		PaymentService: PostgresServiceConfig{
			DB: loadDBConfig("PAYMENT_SERVICE_DB"),
		},
		InventoryService: InventoryServiceConfig{
			DB:    loadDBConfig("INVENTORY_SERVICE_DB"),
			Redis: loadRedisConfig("INVENTORY_SERVICE_REDIS"),
		},
		NotificationService: NotificationServiceConfig{
			Mongo: MongoConfig{
				URI:      getEnv("NOTIFICATION_SERVICE_MONGO_URI", "mongodb://localhost:27017"),
				Database: getEnv("NOTIFICATION_SERVICE_MONGO_DATABASE", "notification_db"),
			},
		},
	}

	return cfg, nil
}

func loadDBConfig(prefix string) PostgresConfig {
	return PostgresConfig{
		Host:     getEnv(prefix+"_HOST", "localhost"),
		Port:     getEnv(prefix+"_PORT", "5432"),
		User:     getEnv(prefix+"_USER", "postgres"),
		Password: getEnv(prefix+"_PASSWORD", "password"),
		DBName:   getEnv(prefix+"_DBNAME", ""),
		SSLMode:  getEnv(prefix+"_SSLMODE", "disable"),
	}
}

func loadRedisConfig(prefix string) RedisConfig {
	dbStr := getEnv(prefix+"_DB", "0")
	db, _ := strconv.Atoi(dbStr)

	return RedisConfig{
		Addr:     getEnv(prefix+"_ADDR", "localhost:6379"),
		Password: getEnv(prefix+"_PASSWORD", ""),
		DB:       db,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
