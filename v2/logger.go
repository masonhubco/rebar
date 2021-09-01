package rebar

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LoggerKey    = "rebarLogger"
	RequestIDKey = "requestID"
)

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	With(fields ...zap.Field) *zap.Logger
}

func NewStandardLogger() (Logger, error) {
	config := zap.NewProductionConfig()
	config.Encoding = "console"
	config.EncoderConfig.EncodeLevel = func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		lvl := l.CapitalString()
		if lvl == "WARN" {
			lvl = "WARNING"
		}
		enc.AppendString(lvl + ":")
	}
	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return config.Build()
}

func LoggerFrom(c *gin.Context) Logger {
	if maybeALogger, exists := c.Get(LoggerKey); exists {
		if logger, ok := maybeALogger.(Logger); ok {
			return logger
		}
	}
	defaultLogger, _ := NewStandardLogger()
	return defaultLogger
}

func RequestIDFrom(c *gin.Context) string {
	return c.GetString(RequestIDKey)
}
