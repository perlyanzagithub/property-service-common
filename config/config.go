package config

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"
	"golang.org/x/net/context"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

// PostgreSQLConfig struct to define the PostgreSQL connection parameters
type PostgreSQLConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	Schema   string // Added schema for flexibility
}

// RedisConfig struct to define the Redis connection parameters
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// KafkaConfig struct to define Kafka consumer parameters
type KafkaConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

// Global variables for PostgreSQL, Redis, and Kafka
var (
	DB          *gorm.DB
	RedisClient *redis.Client
	KafkaReader *kafka.Reader
)

// InitializePostgreSQL initializes the PostgreSQL connection using GORM
func InitializePostgreSQL(config PostgreSQLConfig) error {
	// Build the connection string for PostgreSQL
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s search_path=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode, config.Schema)

	// Open a connection using GORM
	var err error
	DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("Unable to connect to PostgreSQL: %v\n", err)
		return err
	}
	log.Println("PostgreSQL connection established successfully!")
	return nil
}

// InitializeRedis initializes the Redis connection
func InitializeRedis(config RedisConfig) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	// Test Redis connection
	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Unable to connect to Redis: %v\n", err)
		return err
	}
	log.Println("Redis connection established successfully!")
	return nil
}

// InitializeKafka initializes the Kafka consumer
func InitializeKafka(config KafkaConfig) error {
	KafkaReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  config.Brokers,
		Topic:    config.Topic,
		GroupID:  config.GroupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	log.Println("Kafka consumer initialized successfully!")
	return nil
}

// ClosePostgreSQL gracefully closes the PostgreSQL database connection
func ClosePostgreSQL() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("Error getting raw database object: %v", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		log.Printf("Error closing PostgreSQL connection: %v", err)
	}
	log.Println("PostgreSQL connection closed.")
}

// CloseRedis gracefully closes the Redis connection
func CloseRedis() {
	if RedisClient != nil {
		RedisClient.Close()
		log.Println("Redis connection closed.")
	}
}

// CloseKafka gracefully closes the Kafka consumer
func CloseKafka() {
	if KafkaReader != nil {
		KafkaReader.Close()
		log.Println("Kafka consumer closed.")
	}
}
