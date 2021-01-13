package app

import (
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gorilla/websocket"
	"github.com/masonhubco/rebar/middleware"
	"github.com/masonhubco/rebar/samples/graphql/api"
	"github.com/masonhubco/rebar/samples/graphql/graph/generated"
	"github.com/masonhubco/rebar/service"
	"github.com/rs/cors"
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

		app.Router.HandleFunc("/query", buildGraphQLHandler())

	}
	return app
}

type handlerInterface interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

func buildGraphQLHandler() func(w http.ResponseWriter, r *http.Request) {
	c := cors.New(cors.Options{
		/// is this going to work?
		AllowedOrigins: []string{"https://*.masonhub.co"},
		// AllowedOrigins:   []string{"*"},
		// AllowCredentials: true,
		AllowedHeaders: []string{"CENSUS_TOKEN", "Content-Type", "Accept-Language"},
	})

	/// Are iframes different headers? or host headers
	gqlSrv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: &generated.Resolver{}}))
	gqlSrv.AddTransport(transport.POST{})    // is that needed?
	gqlSrv.AddTransport(transport.Websocket{ // is that needed?
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})
	c.Handler(gqlSrv)
	return func(w http.ResponseWriter, r *http.Request) {
		c.HandlerFunc(w, r)
	}
}
