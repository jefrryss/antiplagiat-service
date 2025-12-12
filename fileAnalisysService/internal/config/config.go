package config

import "os"

type Config struct {
	FileStorageURL string
	S3Endpoint     string
	S3AccessKey    string
	S3SecretKey    string
	S3Bucket       string
	ServerPort     string
}

func LoadConfig() *Config {
	return &Config{
		FileStorageURL: getEnv("FILE_STORAGE_URL", "http://filestorage:8080"),
		S3Endpoint:     getEnv("S3_ENDPOINT", "minio:9000"),
		S3AccessKey:    getEnv("S3_ACCESS_KEY", "minio"),
		S3SecretKey:    getEnv("S3_SECRET_KEY", "minio123"),
		S3Bucket:       getEnv("S3_BUCKET", "reports"),
		ServerPort:     getEnv("SERVER_PORT", "8090"),
	}
}

func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
