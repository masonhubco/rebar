package service_test

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/masonhubco/rebar/service"
	"github.com/stretchr/testify/assert"
)

type testProcessor struct {
	starts  int
	stops   int
	startFn func() error
	stopFn  func() error
}

func (tp *testProcessor) Start(ctx context.Context) error {
	tp.starts += 1
	return tp.startFn()
}

func (tp *testProcessor) Stop(wg *sync.WaitGroup) error {
	tp.stops += 1
	wg.Done()
	return tp.stopFn()
}

func Test_serviceServe(t *testing.T) {
	t.Run("Easy Start", func(t *testing.T) {
		t.Parallel()
		re := service.New(service.Options{})

		c := make(chan os.Signal, 1)
		c <- os.Interrupt
		err := re.Serve(c)
		assert.NoError(t, err)
	})

	t.Run("Bad Start", func(t *testing.T) {
		t.Parallel()
		re := service.New(service.Options{ShutDownWaitInSec: 1, PanicOnProcessorStartFailure: true})
		tp := &testProcessor{
			startFn: func() error { return errors.New("totally bad thing that happened") },
			stopFn:  func() error { return errors.New("why doesn't anything work") },
		}
		re.AddProcessor(tp)

		c := make(chan os.Signal, 1)

		err := re.Serve(c)
		assert.EqualError(t, err, "ERROR: Rebar failed to start one or more attached processors (and the PanicOnProcessorStartFailure setting is true)")
	})

	t.Run("Easy Stop", func(t *testing.T) {
		t.Parallel()
		re := service.New(service.Options{ShutDownWaitInSec: 1, PanicOnProcessorStartFailure: true})
		tp := &testProcessor{
			startFn: func() error { return nil },
			stopFn:  func() error { return errors.New("why doesn't anything work") },
		}
		re.AddProcessor(tp)

		c := make(chan os.Signal, 1)
		c <- os.Interrupt
		err := re.Serve(c)
		assert.NoError(t, err)
	})

}

func Test_serviceNew(t *testing.T) {
	t.Parallel()

	t.Run("Defaults", func(t *testing.T) {
		r := service.New(service.Options{})
		assert.Equal(t, "development", r.Environment)
		assert.Equal(t, "0.0.0.0:3000", r.Server.Addr)
		assert.Equal(t, 30*time.Second, r.ShutdownWait)
		assert.Equal(t, 15*time.Second, r.Server.WriteTimeout)
		assert.Equal(t, 15*time.Second, r.Server.ReadTimeout)
		assert.Equal(t, 60*time.Second, r.Server.IdleTimeout)
		assert.False(t, r.PanicOnProcessorStartFailure)
		assert.NotNil(t, r.Router)
	})

	t.Run("Custom Values", func(t *testing.T) {
		r := service.New(service.Options{
			Environment:                  "test",
			Port:                         "3310",
			WriteTimeoutInSec:            35,
			ReadTimeoutInSec:             30,
			IdleTimeoutInSec:             120,
			ShutDownWaitInSec:            60,
			PanicOnProcessorStartFailure: true,
		})
		assert.Equal(t, "test", r.Environment)
		assert.Equal(t, "0.0.0.0:3310", r.Server.Addr)
		assert.Equal(t, 60*time.Second, r.ShutdownWait)
		assert.Equal(t, 35*time.Second, r.Server.WriteTimeout)
		assert.Equal(t, 30*time.Second, r.Server.ReadTimeout)
		assert.Equal(t, 120*time.Second, r.Server.IdleTimeout)
		assert.True(t, r.PanicOnProcessorStartFailure)
		assert.NotNil(t, r.Router)
	})
}

func Test_Processors(t *testing.T) {
	t.Parallel()

	t.Run("Good Start and Good Stop", func(t *testing.T) {
		r := service.New(service.Options{})
		tp := &testProcessor{
			startFn: func() error { return nil },
			stopFn:  func() error { return nil },
		}
		r.AddProcessor(tp)
		errs := r.StartProcessors()
		assert.ElementsMatch(t, []error{}, errs)

		wg := &sync.WaitGroup{}
		errs = r.StopProcessors(wg)
		wg.Wait()
		assert.ElementsMatch(t, []error{}, errs)
		assert.Equal(t, 1, tp.starts)
		assert.Equal(t, 1, tp.stops)
	})

	t.Run("Bad Start and Bad Stop", func(t *testing.T) {
		r := service.New(service.Options{ShutDownWaitInSec: 1})
		tp := &testProcessor{
			startFn: func() error { return errors.New("totally bad thing that happened") },
			stopFn:  func() error { return errors.New("why doesn't anything work") },
		}
		r.AddProcessor(tp)
		errs := r.StartProcessors()
		assert.ElementsMatch(t, []error{errors.New("totally bad thing that happened")}, errs)

		wg := &sync.WaitGroup{}
		errs = r.StopProcessors(wg)
		wg.Wait()
		assert.ElementsMatch(t, []error{errors.New("why doesn't anything work")}, errs)
		assert.Equal(t, 1, tp.starts)
		assert.Equal(t, 1, tp.stops)
	})

}
