# rebar v2

V2 rebar introduced a few improvements

- Simpler subfolder and package structure
- Replaced Gorilla Mux with Gin
- Better [Logger](./middleware/logger.go) middleware with zap
- Improved graceful shutdown with cancelable context
- Updated time duration based option names

### Getting started

```shell
go get github.com/masonhubco/rebar/v2
```

### Examples

Let's start with a minimum example:

```go
logger, err := rebar.NewStandardLogger()
if err != nil {
	log.Fatal("ERROR:", err)
}
app := rebar.New(rebar.Options{
	Environment: rebar.Development,
	Port:        "3000",
	Logger:      logger,
})
app.Router.Use(middleware.Logger(logger))
app.Router.Use(gin.Recovery())

apiGroup := app.Router.Group("/api")
{
	apiGroup.Use(middleware.BasicJWT("test-system-token"))
	apiGroup.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "up",
			"redis":  "connected",
		})
	})
}
if err := app.Run(); err != nil {
	log.Fatal("ERROR:", err)
}
```

And more examples:

- [standard](./examples/standard): basic and simple setup
- [graphql](./examples/graphql): graphQL query and playground integrated
- [graceful](./examples/graceful): gracefully shut down worker or other long running goroutines

Each example app was created as separate go module. To run an example (ie. `standard`):

```shell
cd examples/standard
make run
```

### Configuration

```go
type Options struct {
	// Environment defaults to development. Possible value could be development,
	// test, staging, integration, sandbox and production. When it's set to
	// development, it activates Gin's debug mode, test triggers test mode, and
	// everything else maps to release mode
	Environment string
	// Port defaults to 3000. It's the port rebar http server will listen to.
	Port string
	// Logger is used throughout rebar for writing logs. It accepts an instance
	// of zap logger.
	Logger Logger
	// WriteTimeout defaults to 15 seconds. It maps to http.Server's WriteTimeout.
	WriteTimeout time.Duration
	// ReadTimeout defaults to 15 seconds. It maps to http.Server's ReadTimeout.
	ReadTimeout time.Duration
	// IdleTimeout defaults to 60 seconds. It maps to http.Server's IdleTimeout.
	IdleTimeout time.Duration
	// ShutDownWait defaults to 30 seconds. It tells the server how long it has
	// to gracefully shutdown
	ShutDownWait time.Duration
	// StopOnProcessorStartFailure will prevent the server from starting if any attached processors fail to start
	StopOnProcessorStartFailure bool
}
```
