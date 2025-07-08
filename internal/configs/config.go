package configs

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort          string
	AppEnv           string
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	DBSSLMode        string
	RedisHost        string
	RedisPort        string
	RedisPassword    string
	RedisDB          string
	JWTSecret        string
	JWTExpire        time.Duration
	JWTRefreshExpire time.Duration
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtExpire, _ := time.ParseDuration(os.Getenv("JWT_EXPIRE"))
	jwtRefreshExpire, _ := time.ParseDuration(os.Getenv("JWT_REFRESH_EXPIRE"))

	return &Config{
		AppPort:          os.Getenv("APP_PORT"),
		AppEnv:           os.Getenv("APP_ENV"),
		DBHost:           os.Getenv("DB_HOST"),
		DBPort:           os.Getenv("DB_PORT"),
		DBUser:           os.Getenv("DB_USER"),
		DBPassword:       os.Getenv("DB_PASSWORD"),
		DBName:           os.Getenv("DB_NAME"),
		DBSSLMode:        os.Getenv("DB_SSLMODE"),
		RedisHost:        os.Getenv("REDIS_HOST"),
		RedisPort:        os.Getenv("REDIS_PORT"),
		RedisPassword:    os.Getenv("REDIS_PASSWORD"),
		RedisDB:          os.Getenv("REDIS_DB"),
		JWTSecret:        os.Getenv("JWT_SECRET"),
		JWTExpire:        jwtExpire,
		JWTRefreshExpire: jwtRefreshExpire,
	}
}
