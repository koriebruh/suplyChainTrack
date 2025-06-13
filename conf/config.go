package conf

import (
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
	Version        string
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
	appPort, _ := strconv.Atoi(GetEnv("APP_PORT", "3000"))
	dbPort, _ := strconv.Atoi(GetEnv("DB_PORT", "3000"))

	log.Printf("MODE %v | %v Using APP_PORT: %d, | DB_PORT: %d",
		GetEnv("APP_ENV", "dev-bg"),
		GetEnv("APP_NAME", "SupplyChainTracker -development"),
		appPort, dbPort)

	return &Config{
		AppConfig: AppConfig{
			Version:        GetEnv("APP_VERSION", ""),
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
