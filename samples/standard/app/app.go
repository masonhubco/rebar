package app

import (
	"github.com/masonhubco/rebar/middleware"
	"github.com/masonhubco/rebar/samples/standard/api"
	"github.com/masonhubco/rebar/service"
)

var app *service.Rebar

func App() *service.Rebar {
	if app == nil {
		app = service.New(service.Options{
			Environment: "development",
			Port:        "3005",
		})

		apiSubRouter := app.Router.PathPrefix("/api").Subrouter()
		auth := middleware.AuthenticationMW{SystemToken: "blah"}
		apiSubRouter.Use(auth.Authenticate)
		apiSubRouter.Use(middleware.Logger)
		apiSubRouter.HandleFunc("/status", api.Status())
	}
	return app
}
