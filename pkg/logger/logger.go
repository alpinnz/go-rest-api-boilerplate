// Package logger provides a thin wrapper around zerolog to standardize structured
// logging across the project, and convenience helpers for common logging patterns.
package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger wraps zerolog.Logger for structured logging
type Logger struct {
	logger zerolog.Logger
}

// Config represents logger configuration
type Config struct {
	Level      string // debug, info, warn, error
	Pretty     bool   // human-readable output for development
	Output     io.Writer
	TimeFormat string
}

// New creates a new structured logger
func New(cfg Config) *Logger {
	// Parse log level
	level := parseLevel(cfg.Level)
	zerolog.SetGlobalLevel(level)

	// Configure output
	var output io.Writer = os.Stdout
	if cfg.Output != nil {
		output = cfg.Output
	}

	// Pretty printing for development
	if cfg.Pretty {
		output = zerolog.ConsoleWriter{
			Out:        output,
			TimeFormat: time.RFC3339,
		}
	}

	// Configure time format
	timeFormat := time.RFC3339
	if cfg.TimeFormat != "" {
		timeFormat = cfg.TimeFormat
	}
	zerolog.TimeFieldFormat = timeFormat

	logger := zerolog.New(output).With().
		Timestamp().
		Caller().
		Logger()

	return &Logger{logger: logger}
}

// Default creates a logger with default configuration
func Default() *Logger {
	return New(Config{
		Level:  "info",
		Pretty: false,
	})
}

// Development creates a logger for development environment
func Development() *Logger {
	return New(Config{
		Level:  "debug",
		Pretty: true,
	})
}

// Production creates a logger for production environment
func Production() *Logger {
	return New(Config{
		Level:  "info",
		Pretty: false,
	})
}

// parseLevel parses string log level
func parseLevel(level string) zerolog.Level {
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}

// Logging methods

func (l *Logger) Debug() *zerolog.Event {
	return l.logger.Debug()
}

func (l *Logger) Info() *zerolog.Event {
	return l.logger.Info()
}

func (l *Logger) Warn() *zerolog.Event {
	return l.logger.Warn()
}

func (l *Logger) Error() *zerolog.Event {
	return l.logger.Error()
}

func (l *Logger) Fatal() *zerolog.Event {
	return l.logger.Fatal()
}

// With fields

func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	event := l.logger.With()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	return &Logger{logger: event.Logger()}
}

func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{logger: l.logger.With().Interface(key, value).Logger()}
}

func (l *Logger) WithError(err error) *Logger {
	return &Logger{logger: l.logger.With().Err(err).Logger()}
}

// Context-aware logging

func (l *Logger) WithRequestID(requestID string) *Logger {
	return l.WithField("request_id", requestID)
}

func (l *Logger) WithUserID(userID string) *Logger {
	return l.WithField("user_id", userID)
}

// Global logger functions for backward compatibility

func Debug(msg string) {
	log.Debug().Msg(msg)
}

func Info(msg string) {
	log.Info().Msg(msg)
}

func Warn(msg string) {
	log.Warn().Msg(msg)
}

func Error(msg string) {
	log.Error().Msg(msg)
}

func Fatal(msg string) {
	log.Fatal().Msg(msg)
}

func WithError(err error) *zerolog.Event {
	return log.Error().Err(err)
}

func WithField(key string, value interface{}) *zerolog.Event {
	return log.Info().Interface(key, value)
}

func WithFields(fields map[string]interface{}) *zerolog.Event {
	event := log.Info()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	return event
}
