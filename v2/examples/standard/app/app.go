package app

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/masonhubco/rebar/v2"
	"github.com/masonhubco/rebar/v2/examples/standard/api"
	"github.com/masonhubco/rebar/v2/examples/standard/models"
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
	return app, func() error {
		// waiting for shutdown signal
		<-ctx.Done()
		// shutdown received, stopping...
		return nil
	}, nil
}
