package rebar

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Options is the set of custom options you'd like to use
// to start up this web server.
type Options struct {
	Environment       string
	Port              string
	Logger            Logger
	WriteTimeoutInSec time.Duration
	ReadTimeoutInSec  time.Duration
	IdleTimeoutInSec  time.Duration
	// ShutDownWaitInSec tells the server how long it has to gracefully shutdown
	ShutDownWaitInSec time.Duration
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
	if o.WriteTimeoutInSec.Seconds() == 0 {
		o.WriteTimeoutInSec = 15 * time.Second
	}
	if o.ReadTimeoutInSec.Seconds() == 0 {
		o.ReadTimeoutInSec = 15 * time.Second
	}
	if o.IdleTimeoutInSec.Seconds() == 0 {
		o.IdleTimeoutInSec = 60 * time.Second
	}
	if o.ShutDownWaitInSec.Seconds() == 0 {
		o.ShutDownWaitInSec = 30 * time.Second
	}
	return o
}

const (
	Development = "development"
	Test        = "test"
	Staging     = "staging"
	Integration = "integration"
	Production  = "production"
)

func (o Options) Mode() string {
	switch o.Environment {
	case Development:
		return gin.DebugMode
	case Test:
		return gin.TestMode
	case Staging, Integration, Production:
		return gin.ReleaseMode
	}
	return gin.ReleaseMode
}
