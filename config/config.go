package config

import (
	"os"
)

type Config struct {
	DatabaseURL  string
	KafkaBrokers string
	KafkaTopic   string
}

func LoadConfig() Config {
	return Config{
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		KafkaBrokers: os.Getenv("KAFKA_BROKERS"),
		KafkaTopic:   os.Getenv("KAFKA_TOPIC"),
	}
}
