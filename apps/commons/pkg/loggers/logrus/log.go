package logrus

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/loggers"
)

func convertLogLevel(lvl loggers.Level) logrus.Level {
	switch lvl {
	case loggers.Info:
		return logrus.ErrorLevel
	case loggers.Warning:
		return logrus.WarnLevel
	default:
		return logrus.InfoLevel
	}
}

type Logger struct {
	serviceName string
	logger      *logrus.Logger
}

func NewLogrusLogger(serviceName string) *Logger {
	l := logrus.New()
	l.SetLevel(logrus.InfoLevel)
	l.SetFormatter(&logrus.JSONFormatter{})

	return &Logger{
		serviceName: serviceName,
		logger:      l,
	}
}

func (l *Logger) Log(ctx context.Context, lvl loggers.Level, message string, attrs map[string]interface{}) {
	fs := logrus.Fields{}
	for k, v := range attrs {
		fs[k] = v
	}
	fs["service.name"] = l.serviceName
	l.logger.WithFields(fs).Log(convertLogLevel(lvl), message)
}
