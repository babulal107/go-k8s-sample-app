package internal

import (
	"log"
	"os"
)

type Config struct {
	AppName        string `json:"appName"`
	AppPort        string `json:"appPort"`
	AppLogLevel    string `json:"appLogLevel"`
	DatabaseConfig `json:"databaseConfig"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func LoadConfig() *Config {

	// Read environment variable
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080" // Default value if not set
	}

	logLevel := os.Getenv("APP_LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info" // Default log-level will be `info` if not set
	}

	dbUserName := os.Getenv("DB_USERNAME")
	if dbUserName == "" {
		log.Printf("DB_USERNAME secret var not set\n")
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		log.Printf("DB_PASSWORD secret var not set\n")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Printf("DB_NAME secret var not set\n")
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		log.Printf("DB_PORT secret var not set\n")
	}

	return &Config{
		AppName:     "go-k8s-sample-service",
		AppPort:     port,
		AppLogLevel: logLevel,
		DatabaseConfig: DatabaseConfig{
			Host:     "http://localhost",
			Port:     dbPort,
			Username: dbUserName,
			Password: dbPassword,
			Name:     dbName,
		},
	}

}
