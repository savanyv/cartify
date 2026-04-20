package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName string
	AppEnv  string
	AppPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	JWTSecret             string
	JWTExpiryHours        int
	JWTRefreshExpiryHours int

	CORSAllowedOrigins string
	APIKey             string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		AppName: loadEnv("APP_NAME"),
		AppEnv:  loadEnv("APP_ENV"),
		AppPort: loadEnv("APP_PORT"),

		DBHost:     loadEnv("DB_HOST"),
		DBPort:     loadEnv("DB_PORT"),
		DBUser:     loadEnv("DB_USER"),
		DBPassword: loadEnv("DB_PASSWORD"),
		DBName:     loadEnv("DB_NAME"),

		JWTSecret:             loadEnv("JWT_SECRET"),
		JWTExpiryHours:        loadEnvInt("JWT_EXPIRY_HOURS"),
		JWTRefreshExpiryHours: loadEnvInt("JWT_REFRESH_EXPIRY_HOURS"),

		CORSAllowedOrigins: loadEnv("CORS_ALLOWED_ORIGINS"),
		APIKey:             loadEnv("API_KEY"),
	}
}

func loadEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic("environment variable " + key + " not set")
	}

	return value
}

func loadEnvInt(key string) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic("environment variable " + key + " not set")
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		panic("environment variable " + key + " must be valid integer")
	}

	return intValue
}

func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}

func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development"
}
