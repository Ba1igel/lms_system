package config

import (
	"fmt"
	"log"
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

	cfg := &Config{
		KeycloakURL:      os.Getenv("KEYCLOAK_URL"),
		KeycloakRealm:    os.Getenv("KEYCLOAK_REALM"),
		KeycloakClientID: os.Getenv("KEYCLOAK_CLIENT_ID"),
		KeycloakSecret:   os.Getenv("KEYCLOAK_SECRET"),
		AppPort:          os.Getenv("APP_PORT"),
	}

	if err := cfg.validate(); err != nil {
		log.Fatalf("invalid config: %v", err)
	}

	return cfg
}

func (c *Config) validate() error {
	required := map[string]string{
		"KEYCLOAK_URL":       c.KeycloakURL,
		"KEYCLOAK_REALM":     c.KeycloakRealm,
		"KEYCLOAK_CLIENT_ID": c.KeycloakClientID,
		"KEYCLOAK_SECRET":    c.KeycloakSecret,
	}

	for key, val := range required {
		if val == "" {
			return fmt.Errorf("required environment variable %q is not set", key)
		}
	}

	return nil
}
