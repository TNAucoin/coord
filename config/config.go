package config

import (
	"log"
	"os"
	"strconv"
)

func GetEnv() string {
	return getEnvironmentVariableValue("ENV")
}

func GetApplicationPort() int {
	port := getEnvironmentVariableValue("APPLICATION_PORT")
	portStr, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalf("Failed to convert port to integer %v", err)
	}
	return portStr
}

func GetDatabaseConnectionURL() string {
	return getEnvironmentVariableValue("DATABASE_CONNECTION_URL")
}

func getEnvironmentVariableValue(key string) string {
	if os.Getenv(key) == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return os.Getenv(key)
}
