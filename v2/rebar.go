package rebar

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// Rebar is the MasonHub Base App
type Rebar struct {
	Environment                 string
	ShutdownWait                time.Duration
	StopOnProcessorStartFailure bool
	Router                      *gin.Engine
	Server                      *http.Server
	ctx                         context.Context
	processors                  []Processor
}

// New creates a new Rebar instance. It does not start it up yet....nope, just creates a new Rebar app
// given your supplied options.
//
// Default Settings:
// - Environment: development
// - Port: 3000
// - ShutdownWait: 30 seconds
// - WriteTimeout: 15 seconds
// - ReadTimeout: 15 seconds
// - IdleTimeout: 60 seconds
func New(opts Options) *Rebar {
	opts = opts.ValuesOrDefaults()
	gin.SetMode(opts.Mode())

	router := gin.New()
	return &Rebar{
		Environment:                 opts.Environment,
		Router:                      router,
		StopOnProcessorStartFailure: opts.StopOnProcessorStartFailure,
		ShutdownWait:                opts.ShutDownWait,
		Server: &http.Server{
			Addr:           fmt.Sprintf("0.0.0.0:%s", opts.Port),
			WriteTimeout:   opts.WriteTimeout,
			ReadTimeout:    opts.ReadTimeout,
			IdleTimeout:    opts.IdleTimeout,
			Handler:        router,
			MaxHeaderBytes: 1 << 20,
		},
	}
}

// Serve starts the rebar server and your app.
// func (r *Rebar) Serve(quit <-chan os.Signal) error {
func (r *Rebar) RunWithContext(ctx context.Context, stop context.CancelFunc) error {
	if errs := r.StartProcessors(); len(errs) > 0 {
		if r.StopOnProcessorStartFailure {
			return errors.New("[rebar] ERROR: rebar failed to start one or more attached processors (and the StopOnProcessorStartFailure setting is true)")
		}
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := r.Server.ListenAndServe(); err != nil {
			log.Println("[rebar]", err)
			if !errors.Is(err, http.ErrServerClosed) {
				stop()
			}
		}
	}()

	// Block until we receive our signal.
	<-ctx.Done()
	log.Println("[rebar] shutting down server...")

	var wg sync.WaitGroup
	if errs := r.StopProcessors(&wg); len(errs) > 0 {
		log.Println("[rebar] ERROR: rebar failed to gracefully shutdown one or more attached processors")
	}
	wg.Wait()

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), r.ShutdownWait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	if err := r.Server.Shutdown(ctx); err != nil {
		return err
	}

	log.Println("[rebar] server exiting")
	return nil
}

func (r *Rebar) Run() error {
	return r.RunWithContext(ContextWithCancel())
}

func ContextWithCancel() (context.Context, context.CancelFunc) {
	ctx, stop := context.WithCancel(context.Background())
	CancelOnSignal(stop)
	return ctx, stop
}

func CancelOnSignal(stop context.CancelFunc) {
	// wait for interrupt signal to gracefully shutdown the server with a timeout
	sig := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-sig
		log.Printf("[rebar] system signal received: %+v", s)
		stop()
	}()
}
