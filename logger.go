package pyrus

import "go.uber.org/zap"

// Logger allows you to pass own logger implementation that the library will use.
// By default logger is turned off. Pass *zap.Logger instance to WithZapLogger or you own with more generic WithLogger.
type Logger interface {
	Error(msg string, err error)
}

type zapLogger struct {
	logger *zap.Logger
}

func (l *zapLogger) Error(msg string, err error) {
	l.logger.Error(msg, zap.Error(err))
}

type noopLogger struct{}

func (l *noopLogger) Error(string, error) {}
