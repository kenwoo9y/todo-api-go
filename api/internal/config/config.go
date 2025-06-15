package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port        int
	DBType      string
	DBHost      string
	DBPort      int
	DBName      string
	DBUser      string
	DBPass      string
	CORSOrigins []string
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

	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		return nil, fmt.Errorf("DB_TYPE is required")
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPortStr := os.Getenv("DB_PORT")
	if dbPortStr == "" {
		switch dbType {
		case "mysql":
			dbPortStr = "3306"
		case "postgres":
			dbPortStr = "5432"
		default:
			return nil, fmt.Errorf("unsupported database type: %s", dbType)
		}
	}

	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		return nil, fmt.Errorf("invalid database port number: %w", err)
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		return nil, fmt.Errorf("DB_NAME is required")
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		return nil, fmt.Errorf("DB_USER is required")
	}

	dbPass := os.Getenv("DB_PASSWORD")
	if dbPass == "" {
		return nil, fmt.Errorf("DB_PASSWORD is required")
	}

	// CORSの設定を取得
	corsOriginsStr := os.Getenv("CORS_ORIGINS")
	if corsOriginsStr == "" {
		return nil, fmt.Errorf("CORS_ORIGINS is required")
	}

	corsOrigins := strings.Split(corsOriginsStr, ",")
	if len(corsOrigins) == 0 {
		return nil, fmt.Errorf("CORS_ORIGINS must contain at least one origin")
	}

	return &Config{
		Port:        port,
		DBType:      dbType,
		DBHost:      dbHost,
		DBPort:      dbPort,
		DBName:      dbName,
		DBUser:      dbUser,
		DBPass:      dbPass,
		CORSOrigins: corsOrigins,
	}, nil
}
