package rebar_test

import (
	"context"
	"errors"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/masonhubco/rebar/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockProcessor struct {
	starts  int
	stops   int
	startFn func() error
	stopFn  func() error
}

func (tp *mockProcessor) Start(ctx context.Context) error {
	tp.starts += 1
	return tp.startFn()
}

func (tp *mockProcessor) Stop(wg *sync.WaitGroup) error {
	tp.stops += 1
	wg.Done()
	return tp.stopFn()
}

func Test_Rebar_RunWithContext(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		givenOptions  rebar.Options
		mockProcessor *mockProcessor
		wantError     error
	}{
		{
			name:          "easy start",
			givenOptions:  rebar.Options{},
			mockProcessor: nil,
			wantError:     nil,
		},
		{
			name: "bad start",
			givenOptions: rebar.Options{
				ShutDownWaitInSec:           time.Second,
				StopOnProcessorStartFailure: true,
			},
			mockProcessor: &mockProcessor{
				startFn: func() error { return errors.New("totally bad thing that happened") },
				stopFn:  func() error { return errors.New("why doesn't anything work") },
			},
			wantError: errors.New("[rebar] ERROR: rebar failed to start one or more attached processors (and the StopOnProcessorStartFailure setting is true)"),
		},
		{
			name: "easy stop",
			givenOptions: rebar.Options{
				ShutDownWaitInSec:           time.Second,
				StopOnProcessorStartFailure: true,
			},
			mockProcessor: &mockProcessor{
				startFn: func() error { return nil },
				stopFn:  func() error { return errors.New("why doesn't anything work") },
			},
			wantError: nil,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r := rebar.New(tc.givenOptions)
			if tc.mockProcessor != nil {
				r.AddProcessor(tc.mockProcessor)
			}

			ctx, stop := context.WithCancel(context.Background())
			stop()

			err := r.RunWithContext(ctx, stop)
			if tc.wantError != nil {
				assert.EqualError(t, err, tc.wantError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_Rebar_New(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		givenOptions      rebar.Options
		wantEnvironment   string
		wantServerAddr    string
		wantShutdownWait  time.Duration
		wantWriteTimeout  time.Duration
		wantReadTimeout   time.Duration
		wantIdleTimeout   time.Duration
		wantStopOnFailure bool
	}{
		{
			name:              "defaults",
			givenOptions:      rebar.Options{},
			wantEnvironment:   "development",
			wantServerAddr:    "0.0.0.0:3000",
			wantShutdownWait:  30 * time.Second,
			wantWriteTimeout:  15 * time.Second,
			wantReadTimeout:   15 * time.Second,
			wantIdleTimeout:   60 * time.Second,
			wantStopOnFailure: false,
		},
		{
			name: "custom values",
			givenOptions: rebar.Options{
				Environment:                 "test",
				Port:                        "3310",
				WriteTimeoutInSec:           35 * time.Second,
				ReadTimeoutInSec:            30 * time.Second,
				IdleTimeoutInSec:            120 * time.Second,
				ShutDownWaitInSec:           60 * time.Second,
				StopOnProcessorStartFailure: true,
			},
			wantEnvironment:   "test",
			wantServerAddr:    "0.0.0.0:3310",
			wantShutdownWait:  60 * time.Second,
			wantWriteTimeout:  35 * time.Second,
			wantReadTimeout:   30 * time.Second,
			wantIdleTimeout:   120 * time.Second,
			wantStopOnFailure: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r := rebar.New(tc.givenOptions)
			assert.Equal(t, tc.wantEnvironment, r.Environment)
			assert.Equal(t, tc.wantServerAddr, r.Server.Addr)
			assert.Equal(t, tc.wantShutdownWait, r.ShutdownWait)
			assert.Equal(t, tc.wantWriteTimeout, r.Server.WriteTimeout)
			assert.Equal(t, tc.wantReadTimeout, r.Server.ReadTimeout)
			assert.Equal(t, tc.wantIdleTimeout, r.Server.IdleTimeout)
			assert.Equal(t, tc.wantStopOnFailure, r.StopOnProcessorStartFailure)
			assert.NotNil(t, r.Router)
		})
	}
}

func Test_Rebar_Run(t *testing.T) {
	r := rebar.New(rebar.Options{})
	go func() {
		// wait for 100 millisecond and then send an interrupt signal
		// to the current process, and it should trigger rebar's Run()
		// to unblock and return
		time.Sleep(100 * time.Millisecond)
		p, err := os.FindProcess(syscall.Getpid())
		require.NoError(t, err)
		err = p.Signal(os.Interrupt)
		require.NoError(t, err)
	}()
	r.Run()
}
