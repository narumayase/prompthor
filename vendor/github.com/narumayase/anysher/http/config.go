package http

import (
	"github.com/rs/zerolog"
	"strings"
)

// Config contains the application configuration
type Config struct {
	LogLevel string
}

// NewConfiguration creates configuration for HTTP implementation
func NewConfiguration(logLevel string) Config {
	setLogLevel(logLevel)
	return Config{
		LogLevel: logLevel,
	}
}

// setLogLevel sets the log level defined
func setLogLevel(logLevel string) {
	levels := map[string]zerolog.Level{
		"debug": zerolog.DebugLevel,
		"info":  zerolog.InfoLevel,
		"warn":  zerolog.WarnLevel,
		"error": zerolog.ErrorLevel,
		"fatal": zerolog.FatalLevel,
		"panic": zerolog.PanicLevel,
	}
	levelEnv := strings.ToLower(logLevel)

	level, ok := levels[levelEnv]
	if !ok {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
}
