package config

import (
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	KeycloakURL      string
	KeycloakRealm    string
	KeycloakClientID string
	KeycloakSecret   string
	AppPort          string
}

func Load() *Config {
	_ = godotenv.Load()
	return &Config{
		KeycloakURL:      os.Getenv("KEYCLOAK_URL"),
		KeycloakRealm:    os.Getenv("KEYCLOAK_REALM"),
		KeycloakClientID: os.Getenv("KEYCLOAK_CLIENT_ID"),
		KeycloakSecret:   os.Getenv("KEYCLOAK_SECRET"),
		AppPort:          os.Getenv("APP_PORT"),
	}
}
