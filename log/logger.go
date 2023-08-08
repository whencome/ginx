package log

import "log"

// global logger
var logger Logger = new(defaultLogger)

// Logger define a log interface
type Logger interface {
    Debugf(format string, args ...interface{})
    Infof(format string, args ...interface{})
    Errorf(format string, args ...interface{})
}

// defaultLogger is a logger that will be used when no customized logger registered
type defaultLogger struct{}

func (l *defaultLogger) Debugf(format string, args ...interface{}) {
    log.Printf(format, args...)
}

func (l *defaultLogger) Infof(format string, args ...interface{}) {
    log.Printf(format, args...)
}

func (l *defaultLogger) Errorf(format string, args ...interface{}) {
    log.Printf(format, args...)
}

// Use register a global logger
func Use(l Logger) {
    if l != nil {
        logger = l
    }
}

func Debugf(format string, args ...interface{}) {
    logger.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
    logger.Infof(format, args...)
}

func Errorf(format string, args ...interface{}) {
    logger.Errorf(format, args...)
}
