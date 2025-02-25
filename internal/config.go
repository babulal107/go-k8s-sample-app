package internal

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	AppName        string         `json:"appName" yaml:"APP_NAME"`
	AppPort        string         `json:"appPort" yaml:"APP_PORT"`
	AppLogLevel    string         `json:"appLogLevel" yaml:"APP_LOG_LEVEL"`
	AppEnv         string         `json:"appEnv" yaml:"APP_ENV"`
	DatabaseConfig DatabaseConfig `json:"databaseConfig" yaml:"databaseConfig"`
}

type DatabaseConfig struct {
	Host     string `json:"host" yaml:"HOST"`
	Port     string `json:"port" yaml:"DB_PORT"`
	Username string `json:"username" yaml:"DB_USERNAME"`
	Password string `json:"password" yaml:"DB_PASSWORD"`
	Name     string `json:"name" yaml:"DB_NAME"`
}

func LoadConfigFromEnv() *Config {

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

func LoadConfig() *Config {

	// Read environment variable
	port := os.Getenv("APP_PORT")
	if port == "" {
		fmt.Println("default port used: 8080")
		port = "8080" // Default value if not set
	}

	logLevel := os.Getenv("APP_LOG_LEVEL")
	if logLevel == "" {
		fmt.Println("default log level used: info")
		logLevel = "info" // Default log-level will be `info` if not set
	}

	//var config Config

	// Read db_config from a secret file
	filePath := "/opt/config/db_config.yaml"
	//filePathAppConfig := "/opt/config/app_config.yaml"
	//
	//// Load App config from app_config secret file
	//appDataConfig := ReadConfigFileFromSecret(filePathAppConfig)
	//// Parse YAML
	//err := yaml.Unmarshal(appDataConfig, &config)
	//if err != nil {
	//	log.Fatalf("Error parsing YAML: %v", err)
	//}
	//log.Printf("config object : %+v\n", config)

	// Load DB config from db_config secret file
	dbConfigs := ReadDBSecretConfig(filePath)

	//config.DatabaseConfig = dbConfigs

	return &Config{
		AppName:        "go-k8s-sample-service",
		AppPort:        port,
		AppLogLevel:    logLevel,
		DatabaseConfig: dbConfigs,
	}
}

// ReadDBSecretConfig read config from a file which is mounted as db-secret
func ReadDBSecretConfig(filePath string) DatabaseConfig {

	data := ReadConfigFileFromSecret(filePath)
	// Parse YAML
	var dbConfig DatabaseConfig
	err := yaml.Unmarshal(data, &dbConfig)
	if err != nil {
		log.Fatalf("Error parsing YAML: %v", err)
	}
	return dbConfig
}

func ReadConfigFileFromSecret(filePath string) []byte {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	return data
}
