package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	AppPort    string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPass     string
	DBName     string
	DBTimeZone string

	KeycloakURL          string
	KeycloakRealm        string
	KeycloakClientID     string
	KeycloakClientSecret string

	MinIOEndpoint  string
	MinIOAccessKey string
	MinIOSecretKey string
	MinIOBucket    string
	MinIOUseSSL    bool
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found, using environment variables")
	}

	cfg := &Config{
		AppPort:    getEnv("APP_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPass:     getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "lms_main"),
		DBTimeZone: getEnv("DB_TIMEZONE", "Asia/Almaty"),

		KeycloakURL:          getEnv("KEYCLOAK_URL", ""),
		KeycloakRealm:        getEnv("KEYCLOAK_REALM", ""),
		KeycloakClientID:     getEnv("KEYCLOAK_CLIENT_ID", ""),
		KeycloakClientSecret: getEnv("KEYCLOAK_CLIENT_SECRET", ""),

		MinIOEndpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinIOAccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinIOSecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
		MinIOBucket:    getEnv("MINIO_BUCKET", "lms-attachments"),
		MinIOUseSSL:    getEnv("MINIO_USE_SSL", "false") == "true",
	}

	if err := cfg.validate(); err != nil {
		logrus.Fatalf("invalid config: %v", err)
	}

	return cfg
}

// validate проверяет обязательные поля при старте.
// Лучше упасть сразу с понятной ошибкой, чем в рантайме с паникой внутри gocloak.
func (c *Config) validate() error {
	required := map[string]string{
		"KEYCLOAK_URL":           c.KeycloakURL,
		"KEYCLOAK_REALM":         c.KeycloakRealm,
		"KEYCLOAK_CLIENT_ID":     c.KeycloakClientID,
		"KEYCLOAK_CLIENT_SECRET": c.KeycloakClientSecret,
		"MINIO_ENDPOINT":         c.MinIOEndpoint,
		"MINIO_ACCESS_KEY":       c.MinIOAccessKey,
		"MINIO_SECRET_KEY":       c.MinIOSecretKey,
		"MINIO_BUCKET":           c.MinIOBucket,
	}

	for key, val := range required {
		if val == "" {
			return fmt.Errorf("required environment variable %q is not set", key)
		}
	}

	return nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
