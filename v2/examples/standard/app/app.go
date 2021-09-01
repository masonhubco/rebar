package app

import (
	"github.com/gin-gonic/gin"
	"github.com/masonhubco/rebar/v2"
	"github.com/masonhubco/rebar/v2/examples/standard/api"
	"github.com/masonhubco/rebar/v2/middleware"
)

func App() *rebar.Rebar {
	logger, _ := rebar.NewStandardLogger()
	app := rebar.New(rebar.Options{
		Environment: rebar.Development,
		Port:        "3005",
		Logger:      logger,
	})
	app.Router.Use(middleware.Logger(logger))
	app.Router.Use(gin.Recovery())

	apiGroup := app.Router.Group("/api")
	{
		apiGroup.Use(middleware.BasicJWT("blah"))
		apiGroup.GET("/status", api.Status)
	}
	return app
}
