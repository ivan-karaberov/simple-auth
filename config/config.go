package config

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"reflect"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
	"github.com/sirupsen/logrus"
)

// Holds the configuration settings for the application.
type Config struct {
	DBHost                    string          `env:"DB_HOST"`                      // Database host
	DBPort                    string          `env:"DB_PORT"`                      // Database port
	DBUser                    string          `env:"DB_USER"`                      // Database user
	DBName                    string          `env:"DB_NAME"`                      // Database name
	DBPassword                string          `env:"DB_PASSWORD"`                  // Database password
	WebhookURL                string          `env:"WEBHOOK_URL"`                  // Webhook URL for notifications
	AccessTokenExpireMinutes  int16           `env:"ACCESS_TOKEN_EXPIRE_MINUTES"`  // Access token expiration time in minutes
	RefreshTokenExpireMinutes int16           `env:"REFRESH_TOKEN_EXPIRE_MINUTES"` // Refresh token expiration time in minutes
	RSAPrivateKey             *rsa.PrivateKey // RSA private key for signing tokens
	RSAPublicKey              *rsa.PublicKey  // RSA public key for verifying tokens
}

// Loads the configuration from environment variables and RSA key files.
func LoadConfig(ctx context.Context, filename string, privateRsaKeyPath string, publicRsaKeyPath string) *Config {
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

	rsaPrivateKey, err := loadPrivateKey(privateRsaKeyPath)
	if err != nil {
		logrus.WithError(err).Fatal("Error load private rsa key")
	}

	rsaPublicKey, err := loadPublicKey(publicRsaKeyPath)
	if err != nil {
		logrus.WithError(err).Fatal("Error load public rsa key")
	}

	cfg.RSAPrivateKey = rsaPrivateKey
	cfg.RSAPublicKey = rsaPublicKey

	return &cfg
}

// Reads and parses the RSA private key from a file.
func loadPrivateKey(privateKeyPath string) (*rsa.PrivateKey, error) {
	privKeyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("error reading private key file: %v", err)
	}

	block, _ := pem.Decode(privKeyData)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the private key")
	}

	privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing PKCS#8 private key: %v", err)
	}

	rsaPrivKey, ok := privKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not of type *rsa.PrivateKey")
	}

	return rsaPrivKey, nil
}

// Reads and parses the RSA public key from a file.
func loadPublicKey(publicKeyPath string) (*rsa.PublicKey, error) {
	pubKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("error reading public key file: %v", err)
	}

	block, _ := pem.Decode(pubKeyData)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the public key")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing PKIX public key: %v", err)
	}

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not of type *rsa.PublicKey")
	}

	return rsaPubKey, nil
}
