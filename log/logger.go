package log

import (
	"fmt"
	"log"
	"os"
)

// LogLevel define log level
type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelError
)

// global logger
var (
	logger   Logger = new(defaultLogger)
	logLevel        = LevelInfo // default log level
)

// Logger define a log interface
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// SetLogLevel set global log level
func SetLogLevel(level LogLevel) {
	logLevel = level
}

// defaultLogger is a logger that will be used when no customized logger registered
type defaultLogger struct{}

func (l *defaultLogger) Debugf(format string, args ...interface{}) {
	if logLevel <= LevelDebug {
		log.Printf("[DEBUG] "+format, args...)
	}
}

func (l *defaultLogger) Infof(format string, args ...interface{}) {
	if logLevel <= LevelInfo {
		log.Printf("[INFO] "+format, args...)
	}
}

func (l *defaultLogger) Errorf(format string, args ...interface{}) {
	if logLevel <= LevelError {
		log.Printf("[ERROR] "+format, args...)
	}
}

// Use register a global logger
func Use(l Logger) {
	if l != nil {
		logger = l
	}
}

func Debugf(format string, args ...interface{}) {
	if logLevel <= LevelDebug {
		logger.Debugf(format, args...)
	}
}

func Infof(format string, args ...interface{}) {
	if logLevel <= LevelInfo {
		logger.Infof(format, args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if logLevel <= LevelError {
		logger.Errorf(format, args...)
	}
}

// Printf print formatted log message to stderr
func Printf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}
