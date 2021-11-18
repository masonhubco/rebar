package app

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/masonhubco/rebar/v2"
	"github.com/masonhubco/rebar/v2/examples/graphql/api"
	"github.com/masonhubco/rebar/v2/examples/graphql/graph"
	"github.com/masonhubco/rebar/v2/examples/graphql/models"
	"github.com/masonhubco/rebar/v2/middleware"
	"github.com/rs/cors"
)

func New(
	ctx context.Context,
	logger rebar.Logger,
	info models.Status,
) (app *rebar.Rebar, shutdown func() error, err error) {
	app = rebar.New(rebar.Options{
		Environment: "development",
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

	app.Router.POST("/query", graphQLHandler(app.Environment))
	if app.Environment == rebar.Development {
		app.Router.GET("/play", playgroundHandler())
	}
	return app, func() error {
		// waiting for shutdown signal
		<-ctx.Done()
		// shutdown received, stopping...
		return nil
	}, nil
}

func graphQLHandler(environment string) gin.HandlerFunc {
	h := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowedHeaders: []string{"Content-Type", "Accept-Language"},
		Debug:          false,
	})
	gqlSrv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	if environment == rebar.Development {
		//To make graphql schema load properly in playground development environment.
		gqlSrv.Use(extension.Introspection{})
	}
	gqlSrv.AddTransport(transport.Options{})
	gqlSrv.AddTransport(transport.GET{})
	gqlSrv.AddTransport(transport.POST{})

	return func(c *gin.Context) {
		h.Handler(gqlSrv).ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("Playground", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
