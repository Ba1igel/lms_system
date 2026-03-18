package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	AppPort          string
	DBHost           string
	DBPort           string
	DBUser           string
	DBPass           string
	DBName           string
	KeycloakURL      string
	KeycloakRealm    string
	KeycloakClientID string
	KeycloakSecret   string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		logrus.Warn("No .env file found, using environment variables")
	}

	return &Config{
		AppPort:          getEnv("APP_PORT", "8080"),
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           getEnv("DB_PORT", "5432"),
		DBUser:           getEnv("DB_USER", "postgres"),
		DBPass:           getEnv("DB_PASSWORD", "postgres"),
		DBName:           getEnv("DB_NAME", "lms_main"),
		KeycloakURL:      getEnv("KEYCLOAK_URL", "http://keycloak:8080"),
		KeycloakRealm:    getEnv("KEYCLOAK_REALM", "master"),
		KeycloakClientID: getEnv("KEYCLOAK_CLIENT_ID", "lms-client"),
		KeycloakSecret:   getEnv("KEYCLOAK_SECRET", "your-secret"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
