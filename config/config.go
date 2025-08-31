package config

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
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

	GatewayAPIUrl  string
	GatewayEnabled bool
}

// Load loads configuration from environment variables or an .env file
func Load() Config {
	setLogLevel()
	// Load .env file if it exists (ignore error if file doesn't exist)
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading .env file: %v", err)
	}

	return Config{
		Port:     getEnv("PORT", "8080"),
		LogLevel: getEnv("LOG_LEVEL", "info"),

		OpenAIKey:  getEnv("OPENAI_API_KEY", ""),
		GroqAPIKey: getEnv("GROQ_API_KEY", ""),
		GroqUrl:    getEnv("GROQ_URL", "https://api.groq.com/openai/v1/responses"),
		ChatModel:  getEnv("CHAT_MODEL", "openai/gpt-oss-20b"),

		GatewayAPIUrl:  getEnv("GATEWAY_API_URL", "http://anyway:9889"),
		GatewayEnabled: getEnvAsBool("GATEWAY_ENABLED", true),
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		log.Debug().Msgf("ENV %s: %v", key, value)
		return value
	}
	return defaultValue
}

// getEnvAsBool gets an environment variable as a boolean or returns a default value
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		log.Debug().Msgf("ENV %s: %v", key, value)
		return strings.ToLower(value) == "true"
	}
	return defaultValue
}

// setLogLevel sets the log level defined in LOG_LEVEL environment variable
func setLogLevel() {
	levels := map[string]zerolog.Level{
		"debug": zerolog.DebugLevel,
		"info":  zerolog.InfoLevel,
		"warn":  zerolog.WarnLevel,
		"error": zerolog.ErrorLevel,
		"fatal": zerolog.FatalLevel,
		"panic": zerolog.PanicLevel,
	}
	levelEnv := strings.ToLower(getEnv("LOG_LEVEL", "info"))

	level, ok := levels[levelEnv]
	if !ok {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
}
