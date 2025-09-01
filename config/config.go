package config

import (
	"github.com/joho/godotenv"
	anysherlog "github.com/narumayase/anysher/log"
	"github.com/rs/zerolog/log"
	"os"
)

// Config contains the application configuration
type Config struct {
	Port       string
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
		Port:       getEnv("PORT", "8080"),
		OpenAIKey:  getEnv("OPENAI_API_KEY", ""),
		GroqAPIKey: getEnv("GROQ_API_KEY", ""),
		GroqUrl:    getEnv("GROQ_URL", "https://api.groq.com/openai/v1/responses"),
		ChatModel:  getEnv("CHAT_MODEL", "openai/gpt-oss-20b"),
	}
	anysherlog.SetLogLevel()

	return config
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
