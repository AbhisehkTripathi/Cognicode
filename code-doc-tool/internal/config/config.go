package config

import (
	"os"
)

type Config struct {
	Port        string
	UploadPath  string
	OutputPath  string
	MaxFileSize int64
}

func New() *Config {
	return &Config{
		Port:        getEnv("PORT", "3000"),
		UploadPath:  getEnv("UPLOAD_PATH", "./uploads"),
		OutputPath:  getEnv("OUTPUT_PATH", "./output"),
		MaxFileSize: 100 * 1024 * 1024, // 100MB
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
