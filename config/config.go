package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	IdgApiUrl         string
	IdgApiToken       string
	ClientSecret      string
	ClientId          string
	GatewayInstanceID string
}

// LoadConfig reads .env and populates the Config struct
func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return &Config{
		IdgApiUrl:         os.Getenv("API_URL"),
		IdgApiToken:       os.Getenv("AUTHORIZATION_TOKEN"),
		ClientSecret:      os.Getenv("CLIENT_SECRET"),
		ClientId:          os.Getenv("CLIENT_ID"),
		GatewayInstanceID: os.Getenv("GATEWAY_INSTANCE_ID"),
	}
}
