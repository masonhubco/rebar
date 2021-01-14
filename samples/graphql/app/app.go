package app

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/masonhubco/rebar/middleware"
	"github.com/masonhubco/rebar/samples/graphql/api"
	"github.com/masonhubco/rebar/samples/graphql/graph"
	"github.com/masonhubco/rebar/service"
	"github.com/rs/cors"
)

var app *service.Rebar

func App() *service.Rebar {
	if app == nil {
		options := service.Options{
			Environment: "development",
			Port:        "3005",
		}
		app = service.New(options)

		apiSubRouter := app.Router.PathPrefix("/api").Subrouter()
		auth := middleware.AuthenticationMW{SystemToken: "blah"}
		apiSubRouter.Use(auth.Authenticate)
		apiSubRouter.Use(middleware.Logger)
		apiSubRouter.HandleFunc("/status", api.Status())

		app.Router.HandleFunc("/query", buildGraphQLHandler(options.Environment))
		if options.Environment == "development" {
			app.Router.HandleFunc("/play", playground.Handler("Playground", "/query"))
		}
	}
	return app
}

type handlerInterface interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

func buildGraphQLHandler(environment string) func(w http.ResponseWriter, r *http.Request) {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete},
		AllowedHeaders: []string{"Content-Type", "Accept-Language"},
		Debug:          false,
	})
	var gqlSrv *handler.Server
	if environment == "development" {
		//GraphQL playground schema does not load with `New` but works with `NewDefaultServer`
		gqlSrv = handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	} else {
		gqlSrv = handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
		gqlSrv.AddTransport(transport.Options{})
		gqlSrv.AddTransport(transport.GET{})
		gqlSrv.AddTransport(transport.POST{})
		gqlSrv.AddTransport(transport.MultipartForm{})
	}

	return func(w http.ResponseWriter, r *http.Request) {
		c.Handler(gqlSrv).ServeHTTP(w, r)
	}
}
