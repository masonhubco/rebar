package rebar

import (
	"context"
	"log"
	"sync"
)

// Processor interface defines the necessary functions to start and gracefully stop
// any sub process attached to the Rebar instance
type Processor interface {
	Start(ctx context.Context) (err error)
	Stop(wg *sync.WaitGroup) (err error)
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
			log.Printf("[rebar] ERROR: unable to stop processor: %s", err)
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
			log.Printf("[rebar] ERROR: unable to start processor: %s", err)
			errs = append(errs, err)
		}
	}
	return
}
