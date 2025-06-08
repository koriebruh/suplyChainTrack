package conf

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Config struct {
	AppConfig      AppConfig
	DatabaseConfig DatabaseConfig
}

type AppConfig struct {
	Environment    string
	Port           int
	Name           string
	AllowedOrigins []string
	ApiKey         string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

var (
	configLoaded bool
	configMutex  sync.Once
)

func LoadConfig() *Config {
	configMutex.Do(func() {
		if os.Getenv("APP_ENV") != "production" {
			if err := godotenv.Load(); err != nil {
				log.Printf("Warning: .env file not found: %v", err)
			} else {
				log.Println("Loaded .env file")
			}
		}
	})

	appPort, _ := strconv.Atoi(GetEnv("APP_PORT", "3000"))
	dbPort, _ := strconv.Atoi(GetEnv("DB_PORT", "3000"))

	log.Printf("%v Using APP_PORT: %d, DB_PORT: %d", GetEnv("APP_NAME", "SupplyChainTracker -development"), appPort, dbPort)

	return &Config{
		AppConfig: AppConfig{
			Environment:    GetEnv("APP_ENV", "development"),
			Port:           appPort,
			Name:           GetEnv("APP_NAME", "SupplyChainTracker -development"),
			AllowedOrigins: strings.Split(GetEnv("ALLOWED_ORIGINS", "*"), ","),
			ApiKey:         GetEnv("API_KEY", ""),
		},
		DatabaseConfig: DatabaseConfig{
			Host:     GetEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			Username: GetEnv("DB_USER", "root"),
			Password: GetEnv("DB_PASS", ""),
			Database: GetEnv("DB_NAME", "mydb"),
		},
	}
}

func GetEnv(key, fallback string) string {
	// Cek environment variable, jika tidak ada, gunakan fallback(nilai cadangan)
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// IsProduction untuk check environment
func IsProduction() bool {
	return GetEnv("APP_ENV", "development") == "production"
}
