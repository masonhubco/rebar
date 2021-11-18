package app

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/masonhubco/rebar/v2"
	"github.com/masonhubco/rebar/v2/examples/graceful/api"
	"github.com/masonhubco/rebar/v2/examples/graceful/models"
	"github.com/masonhubco/rebar/v2/middleware"
)

func New(
	ctx context.Context,
	logger rebar.Logger,
	info models.Status,
) (app *rebar.Rebar, shutdown func() error, err error) {
	app = rebar.New(rebar.Options{
		Environment: rebar.Development,
		Port:        "3005",
		Logger:      logger,
	})
	app.Router.Use(middleware.Logger(logger))
	app.Router.Use(gin.Recovery())

	apiGroup := app.Router.Group("/api")
	{
		apiGroup.Use(middleware.BasicJWT("blah"))
		apiGroup.GET("/status", api.Status(info))
	}

	// create a worker and start it in a goroutine
	wrkr := newWorker(logger)
	go wrkr.start()

	return app, func() error {
		// waiting for shutdown signal
		<-ctx.Done()

		logger.Info("graceful shutdown signal received, stopping worker...")
		// shutdown received, stopping worker
		wrkr.stop()
		logger.Info("worker gracefully stopped")
		return nil
	}, nil
}
