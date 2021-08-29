package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/masonhubco/rebar/v2"
	"go.uber.org/zap"
)

// LoggerConfig defines the config for Logger middleware.
type LoggerConfig struct {
	// RequestIDField is the header field name of the request ID
	RequestIDField string

	// SkipPaths is a url path array which logs are not written.
	// Optional.
	SkipPaths []string
}

const RequestIDField = "X-Request-ID"

func Logger(logger rebar.Logger) gin.HandlerFunc {
	return LoggerWithConfig(logger, LoggerConfig{
		RequestIDField: RequestIDField,
	})
}

func LoggerWithConfig(logger rebar.Logger, conf LoggerConfig) gin.HandlerFunc {
	notlogged := conf.SkipPaths

	var skip map[string]struct{}
	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		reqID := c.Request.Header.Get(conf.RequestIDField)
		if reqID == "" {
			reqID = uuid.Must(uuid.NewV4()).String()
		}
		c.Header(conf.RequestIDField, reqID)
		c.Set(rebar.RequestIDKey, reqID)
		c.Set(rebar.LoggerKey, logger.With(zap.String("request_id", reqID)))

		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			// Stop timer
			latency := time.Since(start)
			if latency > time.Minute {
				latency = latency.Truncate(time.Second)
			}

			if raw != "" {
				path = path + "?" + raw
			}

			contentType := c.ContentType()
			if contentType == "" {
				contentType = c.Writer.Header().Get("Content-Type")
			}

			msg := "[rebar] " + path
			fields := []zap.Field{
				zap.String(conf.RequestIDField, reqID),
				zap.String("content_type", contentType),
				zap.Int("body_bytes", c.Writer.Size()),
				zap.Int("status_code", c.Writer.Status()),
				zap.String("latency", latency.String()),
				zap.String("client_ip", c.ClientIP()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
			}

			if len(c.Errors) > 0 {
				fields = append(fields,
					zap.String("error", c.Errors.ByType(gin.ErrorTypePrivate).String()))
				logger.Error(msg, fields...)
			} else {
				logger.Info(msg, fields...)
			}
		}
	}
}
