package app

import (
	"log"

	"github.com/masonhubco/rebar/v2"
)

type worker struct {
	logger rebar.Logger
	quit   chan bool
}

func newWorker(logger rebar.Logger) *worker {
	return &worker{
		logger: logger,
		quit:   make(chan bool),
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
