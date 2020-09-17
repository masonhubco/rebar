package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

//Rebar is the MasonHub Base App
type Rebar struct {
	Environment                  string
	ShutdownWait                 time.Duration
	PanicOnProcessorStartFailure bool
	Router                       *mux.Router
	Server                       *http.Server
	ctx                          context.Context
	processors                   []Processor
}

// Options is the set of custom options you'd like to use
// to start up this web server.
type Options struct {
	Environment       string
	Port              string
	WriteTimeoutInSec int
	ReadTimeoutInSec  int
	IdleTimeoutInSec  int
	//ShutDownWaitInSec tells the server how long it has to gracefully shutdown
	ShutDownWaitInSec int
	//PanicOnProcessorStartFailure will prevent the server from starting if any attached processors fail to start
	PanicOnProcessorStartFailure bool
}

// Processor interface defines the necessary functions to start and gracefully stop
// any sub process attached to the Rebar instance
type Processor interface {
	Start(ctx context.Context) (err error)
	Stop(wg *sync.WaitGroup) (err error)
}

// New creates a new Rebar instance. It does not start it up yet....nope, just creates a new Rebar app
// given your supplied options.
// Default Settings:
// *Environment*: development
// *Port*: 3000
// *ShutdownWait*: 30 seconds
// *WriteTimeout*: 15 seconds
// *ReadTimeout*: 15 seconds
// *IdleTimeout*: 60 seconds
func New(ro Options) *Rebar {
	router := mux.NewRouter()
	return &Rebar{
		Environment: func(env string) string {
			if env == "" {
				return "development"
			}
			return ro.Environment
		}(ro.Environment),
		Router:                       router,
		PanicOnProcessorStartFailure: ro.PanicOnProcessorStartFailure,
		ShutdownWait:                 defaultTimeout(ro.ShutDownWaitInSec, 30, time.Second),
		Server: &http.Server{
			Addr: func(port string) string {
				if port == "" {
					port = "3000"
				}
				return fmt.Sprintf("0.0.0.0:%s", port)
			}(ro.Port),
			WriteTimeout: defaultTimeout(ro.WriteTimeoutInSec, 15, time.Second),
			ReadTimeout:  defaultTimeout(ro.ReadTimeoutInSec, 15, time.Second),
			IdleTimeout:  defaultTimeout(ro.IdleTimeoutInSec, 60, time.Second),
			Handler:      router, // Pass our instance of gorilla/mux in.
		},
	}
}

// AddProcessor allows you to hang any additional sub processes off of the web server
// It must conform to the processor interface defined above.
// Rebar will attempt to start and gracefully stop any attached process using the Start and Stop functions
func (r *Rebar) AddProcessor(p Processor) {
	r.processors = append(r.processors, p)
}

// StopProcessors stops the attached processors using the Stop() method from the interface.
// It also builds a list of errors and logs them out so you can do something about those errors.
func (r *Rebar) StopProcessors(wg *sync.WaitGroup) (errs []error) {
	for _, p := range r.processors {
		wg.Add(1)
		err := p.Stop(wg)
		if err != nil {
			log.Printf("ERROR (rebar): unable to stop processor: %s", err)
			errs = append(errs, err)
		}
	}
	return
}

// StartProcessors starts the attached processors and builds
// a list of any errors from starting said processors.
// So you can review and then, you know, do something about them.
func (r *Rebar) StartProcessors() (errs []error) {
	ctx := context.Background()
	for _, p := range r.processors {
		err := p.Start(ctx)
		if err != nil {
			log.Printf("ERROR (rebar): unable to start processor: %s", err)
			errs = append(errs, err)
		}
	}
	return
}

//Serve starts the rebar server and your app.
func (r *Rebar) Serve(c <-chan os.Signal) (err error) {
	startErrs := r.StartProcessors()
	if len(startErrs) > 0 && r.PanicOnProcessorStartFailure {
		return errors.New("ERROR: Rebar failed to start one or more attached processors (and the PanicOnProcessorStartFailure setting is true)")
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := r.Server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	// Block until we receive our signal.
	<-c

	wg := &sync.WaitGroup{}
	stopErrs := r.StopProcessors(wg)
	if len(stopErrs) > 0 {
		log.Println("ERROR: Rebar failed to gracefully shutdown one or more attached processors")
	}
	wg.Wait()

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), r.ShutdownWait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	r.Server.Shutdown(ctx)

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("rebar: shutting down")

	return
}

func defaultTimeout(override int, def int, dur time.Duration) time.Duration {
	if override == 0 {
		return dur * time.Duration(def)
	}
	return time.Duration(override) * dur
}
