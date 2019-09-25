package config

import (
	"os"
)

const (
	environment       = "TODO_ENV"
	privateToken      = "TODO_KEY"
	httpPort          = "TODO_HTTP_PORT"
	db                = "TODO_DB"
	slackToken        = "TODO_SLACK_TOKEN"
	slackAlertChannel = "TODO_SLACK_ALERT_CHANNEL"
)

// Config contains application configuration
type Config struct {
	Environment       string
	PrivateToken      string
	HTTPPort          string
	DB                string
	SlackToken        string
	SlackAlertChannel string
}

var config *Config

func getEnvOrDefault(env string, defaultVal string) string {
	e := os.Getenv(env)
	if e == "" {
		return defaultVal
	}
	return e
}

// GetConfiguration , get application configuration based on set environment
func GetConfiguration() (*Config, error) {
	if config != nil {
		return config, nil
	}

	// default configuration
	config := &Config{
		Environment:       getEnvOrDefault(environment, "dev"),
		PrivateToken:      getEnvOrDefault(privateToken, "test"),
		HTTPPort:          getEnvOrDefault(httpPort, "8080"),
		DB:                getEnvOrDefault(db, "postgres://postgres@localhost:5432/sample?sslmode=disable"),
		SlackToken:        getEnvOrDefault(slackToken, ""),
		SlackAlertChannel: getEnvOrDefault(slackAlertChannel, "#todo-alert"),
	}

	return config, nil
}
