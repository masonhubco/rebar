package main

import (
	"github.com/masonhubco/rebar/v2"
	"github.com/masonhubco/rebar/v2/examples/standard/app"
	"github.com/masonhubco/rebar/v2/examples/standard/models"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

/// Start server by running `make run` in v2/examples/standard folder
/// Query the status endpoint by using curl `curl -H 'Accept: application/json' -H "Authorization: Bearer blah" http://localhost:3005/api/status`
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
