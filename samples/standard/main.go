package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/masonhubco/rebar/samples/standard/app"
)

/// Start server by running `go run samples/standard/main.go`
/// Query the status endpoint by using curl `curl -H 'Accept: application/json' -H "Authorization: Bearer blah" http://localhost:3005/api/status`
func main() {
	app := app.App()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	err := app.Serve(c)
	if err != nil {
		panic(err)
	}
}
