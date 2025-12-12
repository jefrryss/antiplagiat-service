package config

import "os"

type Config struct {
	FileStorageURL  string
	AnalysisService string
	ServerPort      string
}

func LoadConfig() *Config {
	return &Config{
		FileStorageURL:  getEnv("FILE_STORAGE_URL", "http://filestorage:8080"),
		AnalysisService: getEnv("ANALYSIS_SERVICE_URL", "http://analysis:8082"),
		ServerPort:      getEnv("SERVER_PORT", "8080"),
	}
}

func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
