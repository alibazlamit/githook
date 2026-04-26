package config

import "os"

type Config struct {
	Port          string
	DatabaseDSN   string
	NATSUrl       string
	WebhookSecret string
}

func Load() *Config {
	return &Config{
		Port:          getEnv("PORT", "8080"),
		DatabaseDSN:   getEnv("DATABASE_DSN", ""),
		NATSUrl:       getEnv("NATS_URL", "nats://localhost:4222"),
		WebhookSecret: getEnv("WEBHOOK_SECRET", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
