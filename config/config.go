package config

import (
	"github.com/joho/godotenv"
	anysherLog "github.com/narumayase/anysher/log"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

// Config contains the application configuration
type Config struct {
	LogLevel string
	Port     string

	OpenAIKey  string
	GroqAPIKey string
	GroqUrl    string
	ChatModel  string
}

// Load loads configuration from environment variables or an .env file
func Load() Config {
	// Load .env file if it exists (ignore error if file doesn't exist)
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading .env file: %v", err)
	}
	config := Config{
		Port:     getEnv("PORT", "8080"),
		LogLevel: getEnv("LOG_LEVEL", "info"),

		OpenAIKey:  getEnv("OPENAI_API_KEY", ""),
		GroqAPIKey: getEnv("GROQ_API_KEY", ""),
		GroqUrl:    getEnv("GROQ_URL", "https://api.groq.com/openai/v1/responses"),
		ChatModel:  getEnv("CHAT_MODEL", "openai/gpt-oss-20b"),
	}
	anysherLog.SetLogLevel()

	return config
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsBool gets an environment variable as a boolean or returns a default value
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return strings.ToLower(value) == "true"
	}
	return defaultValue
}
