package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configurations
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Kafka    KafkaConfig
	LogLevel string
}

// ServerConfig holds HTTP server settings
type ServerConfig struct {
	Port string
}

// DatabaseConfig holds PostgreSQL settings
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// RedisConfig holds Redis settings
type RedisConfig struct {
	Host string
	Port string
}

// KafkaConfig holds Kafka settings
type KafkaConfig struct {
	Brokers string
}

// Load func reads configurations from environment variables
// it attempts to load .env file but doesnt fail if not found
func Load() *Config {
	godotenv.Load()

	return &Config{
		Server: ServerConfig{
			Port: os.Getenv("SERVER_PORT"),
		},
		Database: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
		},
		Redis: RedisConfig{
			Host: os.Getenv("REDIS_HOST"),
			Port: os.Getenv("REDIS_PORT"),
		},
		Kafka: KafkaConfig{
			Brokers: os.Getenv("KAFKA_BROKERS"),
		},
		LogLevel: os.Getenv("LOG_LEVEL"),
	}
}

// DSN returns PostgreSQL connection string
func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host, c.Database.Port, c.Database.User, c.Database.Password, c.Database.Name, c.Database.SSLMode)
}
