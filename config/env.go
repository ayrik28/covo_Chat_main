package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken     string
	DeepSeekToken     string
	MaxRequestsPerDay int
	CooldownSeconds   int
	// MySQL Config
	MySQLHost     string
	MySQLPort     string
	MySQLUser     string
	MySQLPassword string
	MySQLDatabase string
}

var AppConfig *Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("فایل .env یافت نشد، از متغیرهای محیطی استفاده می‌شود")
	}

	AppConfig = &Config{
		// TelegramToken:     getEnv("TELEGRAM_TOKEN", ""),
		// DeepSeekToken:     getEnv("DEEPSEEK_TOKEN", ""),
		MaxRequestsPerDay: getEnvAsInt("MAX_REQUESTS_PER_DAY", 5),
		CooldownSeconds:   getEnvAsInt("COOLDOWN_SECONDS", 10),
		// MySQL Config
		MySQLHost:     getEnv("MYSQL_HOST", "localhost"),
		MySQLPort:     getEnv("MYSQL_PORT", "3306"),
		MySQLUser:     getEnv("MYSQL_USER", ""),
		MySQLPassword: getEnv("MYSQL_PASSWORD", ""),
		MySQLDatabase: getEnv("MYSQL_DATABASE", "myappdb"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
