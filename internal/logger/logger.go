package logger

import (
	"io"
	"os"
	"time"

	"github.com/jaimesHub/golang-todo-app/internal/config"
	"github.com/sirupsen/logrus"
)

// Logger is a wrapper around logrus.Logger
type Logger struct {
	*logrus.Logger
}

// NewLogger creates a new logger
func NewLogger(cfg config.LoggingConfig) (*Logger, error) {
	logger := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}
	logger.SetLevel(level)

	// Set formatter
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	// Set output
	if cfg.File != "" {
		file, err := os.OpenFile(cfg.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}

		// Write to both file and stdout
		logger.SetOutput(io.MultiWriter(os.Stdout, file))
	} else {
		// Write to stdout only
		logger.SetOutput(os.Stdout)
	}

	return &Logger{logger}, nil
}

// WithField adds a field to the log entry
func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.Logger.WithField(key, value)
}

// WithFields adds multiple fields to the log entry
func (l *Logger) WithFields(fields map[string]interface{}) *logrus.Entry {
	return l.Logger.WithFields(logrus.Fields(fields))
}

// WithError adds an error to the log entry
func (l *Logger) WithError(err error) *logrus.Entry {
	return l.Logger.WithError(err)
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		l.WithFields(fields[0]).Debug(msg)
	} else {
		l.Logger.Debug(msg)
	}
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		l.WithFields(fields[0]).Info(msg)
	} else {
		l.Logger.Info(msg)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		l.WithFields(fields[0]).Warn(msg)
	} else {
		l.Logger.Warn(msg)
	}
}

// Error logs an error message
func (l *Logger) Error(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		l.WithFields(fields[0]).Error(msg)
	} else {
		l.Logger.Error(msg)
	}
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		l.WithFields(fields[0]).Fatal(msg)
	} else {
		l.Logger.Fatal(msg)
	}
}
