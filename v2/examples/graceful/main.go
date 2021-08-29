package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/masonhubco/rebar/v2"
	"github.com/masonhubco/rebar/v2/middleware"
)

// this example demonstrates gracefully shutdown long running goroutines such as worker,
// calling worker's stop method when shutdown signal is received and notified by cancel
// context

func main() {
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

	// create a worker and start it in a goroutine
	wrkr := newWorker()
	go wrkr.start()

	// create a cancelable context and use the context to notify worker for shutdown
	ctx, stop := rebar.ContextWithCancel()
	go func() {
		// waiting for shutdown signal
		<-ctx.Done()
		// shutdown received, stopping worker
		wrkr.stop()
	}()

	// starting rebar app with the cancelable context
	if err := app.RunWithContext(ctx, stop); err != nil {
		log.Fatal("ERROR:", err)
	}
}

type worker struct {
	quit chan bool
}

func newWorker() *worker {
	return &worker{
		quit: make(chan bool),
	}
}

func (w *worker) start() {
	log.Println("[worker] starting...")
	<-w.quit
	log.Println("[worker] stopped")
}

func (w *worker) stop() {
	log.Println("[worker] stopping...")
	w.quit <- true
}
