package config

import "os"

type Config struct {
	AWSEndpoint    string
	AWSRegion      string
	ProcessedQueue string
	Port           string
}

func Load() *Config {
	return &Config{
		AWSEndpoint:    getEnv("AWS_ENDPOINT_URL", "http://localhost:4566"),
		AWSRegion:      getEnv("AWS_REGION", "us-east-1"),
		ProcessedQueue: getEnv("PROCESSED_EVENTS_QUEUE", "processed-events"),
		Port:           getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}