package main

import (
	"log"

	"github.com/masonhubco/rebar/v2/examples/graphql/app"
)

/// Start server by running `make run` in v2/examples/graphql folder
/// Query the status endpoint by using curl `curl -H 'Accept: application/json' -H "Authorization: Bearer blah" http://localhost:3005/api/status`
func main() {
	app := app.App()
	if err := app.Run(); err != nil {
		log.Fatal("ERROR:", err)
	}
}
