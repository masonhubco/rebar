package main

import (
	"github.com/masonhubco/rebar/v2"
	"github.com/masonhubco/rebar/v2/examples/graceful/app"
	"github.com/masonhubco/rebar/v2/examples/graceful/models"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// this example demonstrates gracefully shutdown long running goroutines such as worker,
// calling worker's stop method when shutdown signal is received and notified by cancel
// context

func main() {
	logger, _ := rebar.NewStandardLogger()
	ctx, stop := rebar.ContextWithCancel()
	status := models.NewStatus(gitCommit, buildTime)
	app, shutdown, err := app.New(ctx, logger, status)
	if err != nil {
		logger.Error("failed creating example app", zap.Error(err))
		return
	}
	g := new(errgroup.Group)
	g.Go(func() error {
		// starting rebar app with the cancelable context
		return app.RunWithContext(ctx, stop)
	})
	g.Go(func() error {
		return shutdown()
	})
	if err := g.Wait(); err != nil {
		logger.Error("failed running or shutting down example app", zap.Error(err))
		return
	}
}

var (
	gitCommit string
	buildTime string
)
