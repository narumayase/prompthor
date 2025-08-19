package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

// Config contains the application configuration
type Config struct {
	Port       string
	OpenAIKey  string
	GroqAPIKey string
	ChatModel  string
}

// Load loads configuration from environment variables or an .env file
func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
	return &Config{
		Port:       getEnv("PORT", "8080"),
		OpenAIKey:  getEnv("OPENAI_API_KEY", ""),
		GroqAPIKey: getEnv("GROQ_API_KEY", ""),
		ChatModel:  getEnv("CHAT_MODEL", "openai/gpt-oss-20b"),
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
