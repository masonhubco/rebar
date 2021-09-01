package rebar

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Options is the set of custom options you'd like to use
// to start up this web server.
type Options struct {
	// Environment defaults to development. Possible value could be development,
	// test, staging, integration, sandbox and production. When it's set to
	// development, it activates Gin's debug mode, test triggers test mode, and
	// everything else maps to release mode
	Environment string
	// Port defaults to 3000. It's the port rebar http server will listen to.
	Port string
	// Logger is used throughout rebar for writing logs. It accepts an instance
	// of zap logger.
	Logger Logger
	// WriteTimeout defaults to 15 seconds. It maps to http.Server's WriteTimeout.
	WriteTimeout time.Duration
	// ReadTimeout defaults to 15 seconds. It maps to http.Server's ReadTimeout.
	ReadTimeout time.Duration
	// IdleTimeout defaults to 60 seconds. It maps to http.Server's IdleTimeout.
	IdleTimeout time.Duration
	// ShutDownWait defaults to 30 seconds. It tells the server how long it has
	// to gracefully shutdown
	ShutDownWait time.Duration
	// StopOnProcessorStartFailure will prevent the server from starting if any attached processors fail to start
	StopOnProcessorStartFailure bool
}

func (o Options) ValuesOrDefaults() Options {
	if o.Environment == "" {
		o.Environment = "development"
	}
	if o.Port == "" {
		o.Port = "3000"
	}
	if o.Logger == nil {
		o.Logger, _ = NewStandardLogger()
	}
	if o.WriteTimeout == 0 {
		o.WriteTimeout = 15 * time.Second
	}
	if o.ReadTimeout == 0 {
		o.ReadTimeout = 15 * time.Second
	}
	if o.IdleTimeout == 0 {
		o.IdleTimeout = 60 * time.Second
	}
	if o.ShutDownWait == 0 {
		o.ShutDownWait = 30 * time.Second
	}
	return o
}

const (
	Development = "development"
	Test        = "test"
	Staging     = "staging"
	Sandbox     = "sandbox"
	Integration = "integration"
	Production  = "production"
)

func (o Options) Mode() string {
	switch strings.ToLower(o.Environment) {
	case Development:
		return gin.DebugMode
	case Test:
		return gin.TestMode
	case Staging, Integration, Sandbox, Production:
		return gin.ReleaseMode
	}
	return gin.ReleaseMode
}
