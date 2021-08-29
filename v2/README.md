# rebar v2

V2 rebar introduced a few improvements

- Simpler subfolder and package structure
- Replaced Gorilla Mux with Gin
- Better [Logger](./middleware/logger.go) middleware with zap
- Improved graceful shutdown with cancelable context

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
