# standard

Example for basic and simple setup. Run this example with `make run`.

Nested packages in this example are not required for small apps. Subpackages help organize code
in better structure for larger apps and services.

If it's only a small app, putting everything in the main package (and in separate files if needed)
like the following example would be sufficient:

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

See another example [graceful](../graceful) for small app, and [graphql](../graphql) for slightly larger app.
