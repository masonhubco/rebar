package app

import (
	"github.com/masonhubco/rebar/middleware"
	"github.com/masonhubco/rebar/samples/graphql/api"
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

		gqlSubRouter := app.Router.PathPrefix("/graphql").Subrouter()

		// srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
		gqlHandler := handler.New(gqlapi.NewExecutableSchema(gqlapi.Config{Resolvers: &gqlapi.Resolver{}}))

		app.Router.HandleFunc("/query", HandlerWrap(gqlHandler))

	}
	return app
}

type handlerInterface interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

func (c *Cors) HandlerWrap(hi handlerInterface) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				hi.ServeHTTP(w, r)
	})
}