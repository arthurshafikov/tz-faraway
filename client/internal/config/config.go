package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppConfig
}

type AppConfig struct {
	ServerAddress string
}

func NewConfig(envFileLocation string) (*Config, error) {
	if err := godotenv.Load(envFileLocation); err != nil {
		return nil, err
	}

	return &Config{
		AppConfig: AppConfig{
			ServerAddress: os.Getenv("SERVER_ADDRESS"),
		},
	}, nil
}
