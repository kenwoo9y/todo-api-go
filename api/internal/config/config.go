package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port int
}

func New() (*Config, error) {
	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080"
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid port number: %w", err)
	}

	return &Config{
		Port: port,
	}, nil
}
