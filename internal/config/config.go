package config

import "os"

type Config struct {
	NATSUrl        string
	DatabaseURL    string
	HTTPPort       string
	JaegerEndpoint string
}

func Load() Config {
	return Config{
		NATSUrl:        getEnv("NATS_URL", "nats://localhost:4222"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://postgres:warden@localhost:5432/warden_engine_db"),
		HTTPPort:       getEnv("HTTP_PORT", ":8081"),
		JaegerEndpoint: getEnv("JAEGER_ENDPOINT", "localhost:4318"),
	}
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
