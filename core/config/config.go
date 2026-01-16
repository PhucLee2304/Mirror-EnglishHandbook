package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string
	AppMode string

	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string

	RedisHost      string
	RedisPort      string
	RedisRatelimit string

	DataJsonPath string
}

func LoadConfig() (*Config, error) {
	// if envPath != "" {
	// 	if err := godotenv.Load(envPath); err != nil {
	// 		log.Printf("Error loading .env file: %v", err)
	// 	}
	// } else {
	// 	if err := godotenv.Load(); err != nil {
	// 		log.Printf("Error loading .env file: %v", err)
	// 	}
	// }
	_ = godotenv.Load()

	cfg := &Config{
		AppPort: getEnv("APP_PORT", "8000"),
		AppMode: getEnv("APP_MODE", "development"),

		DBUser:     getEnv("POSTGRES_USER", "user"),
		DBPassword: getEnv("POSTGRES_PASSWORD", "password"),
		DBName:     getEnv("POSTGRES_DB", "dictionary"),
		DBHost:     getEnv("POSTGRES_HOST", "localhost"),
		DBPort:     getEnv("POSTGRES_PORT", "5432"),

		RedisHost:      getEnv("REDIS_HOST", "localhost"),
		RedisPort:      getEnv("REDIS_PORT", "6379"),
		RedisRatelimit: getEnv("REDIS_RATE_LIMIT", "1000"),

		DataJsonPath: getEnv("DATA_JSON_PATH", ""),
	}

	return cfg, nil
}

func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}

func (c *Config) RedisAddress() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
