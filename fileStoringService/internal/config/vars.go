package config

import (
	"os"
)

type Config struct {
	PostgresDBUrl string
	S3Endpoint    string
	S3AccessKey   string
	S3SecretKey   string
	S3Bucket      string
	ServerPort    string
}

var AppConfig *Config

func LoadConfig() *Config {
	AppConfig = &Config{
		PostgresDBUrl: getEnv("DATABASE_URL", "postgres://filestorage:filestorage@localhost:5432/filestorage?sslmode=disable"),
		S3Endpoint:    getEnv("S3_ENDPOINT", "localhost:9000"),
		S3AccessKey:   getEnv("S3_ACCESS_KEY", "minio"),
		S3SecretKey:   getEnv("S3_SECRET_KEY", "minio123"),
		S3Bucket:      getEnv("S3_BUCKET", "files"),
		ServerPort:    getEnv("SERVER_PORT", "8080"),
	}
	return AppConfig
}

func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
