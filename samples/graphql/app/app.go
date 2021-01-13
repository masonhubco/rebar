package app

import (
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
	"github.com/masonhubco/rebar/middleware"
	"github.com/masonhubco/rebar/samples/graphql/api"
	"github.com/masonhubco/rebar/samples/graphql/graph"
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
		app.Router.HandleFunc("/play", playground.Handler("Test", "/query"))

	}
	return app
}

type handlerInterface interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

func buildGraphQLHandler() func(w http.ResponseWriter, r *http.Request) {
	c := cors.New(cors.Options{
		/// is this going to work?
		// AllowedOrigins: []string{"https://*.masonhub.co"},
		AllowedOrigins: []string{"*"},
		// AllowCredentials: true,
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete},
		AllowedHeaders: []string{"Content-Type", "Accept-Language"},
		Debug:          true,
	})

	/// Are iframes different headers? or host headers
	gqlSrv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	gqlSrv.AddTransport(transport.Options{})
	gqlSrv.AddTransport(transport.GET{})
	gqlSrv.AddTransport(transport.POST{})
	gqlSrv.AddTransport(transport.MultipartForm{})
	gqlSrv.AddTransport(transport.Websocket{ // is that needed?
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})

	return func(w http.ResponseWriter, r *http.Request) {
		c.Handler(gqlSrv).ServeHTTP(w, r)
	}
}
