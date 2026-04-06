package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_PASSWORD string
	JWT_SECRET  string
	BD_HOST     string
}

func initConfig() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("no .env file found (ok for prod):", err)
	}
	return &Config{
		DB_PASSWORD: os.Getenv("PG_USER"),
		JWT_SECRET:  os.Getenv(""),
		BD_HOST:     os.Getenv(""),
	}
}
