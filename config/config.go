package config

import (
	"context"
	"reflect"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
	"github.com/sirupsen/logrus"
)

type Config struct {
	DBHost     string `env:"DB_HOST"`
	DBPort     string `env:"DB_PORT"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	WebhookURL string `env:"WEBHOOK_URL"`
}

func LoadConfig(ctx context.Context, filename string) *Config {
	if err := godotenv.Load(filename); err != nil {
		logrus.WithError(err).Fatal("Error loading .env file")
	}

	var cfg Config
	if err := envconfig.Process(ctx, &cfg); err != nil {
		logrus.WithError(err).Error("Error parse .env file")
	}

	v := reflect.ValueOf(cfg)
	for field_idx := range v.NumField() {
		if v.Field(field_idx).IsZero() {
			logrus.Fatalf("Missing required environment variable: %s", v.Type().Field(field_idx).Tag.Get("env"))
		}
	}

	return &cfg
}
