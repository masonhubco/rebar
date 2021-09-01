package rebar_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/masonhubco/rebar/v2"
	"github.com/stretchr/testify/assert"
)

func Test_Processors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		givenOptions  rebar.Options
		mockProcessor mockProcessor
		wantStartErrs []error
		wantStopErrs  []error
	}{
		{
			name:         "good start and good stop",
			givenOptions: rebar.Options{},
			mockProcessor: mockProcessor{
				startFn: func() error { return nil },
				stopFn:  func() error { return nil },
			},
			wantStartErrs: []error{},
			wantStopErrs:  []error{},
		},
		{
			name: "bad start and bad stop",
			givenOptions: rebar.Options{
				ShutDownWait: 1 * time.Second,
			},
			mockProcessor: mockProcessor{
				startFn: func() error { return errors.New("totally bad thing that happened") },
				stopFn:  func() error { return errors.New("why doesn't anything work") },
			},
			wantStartErrs: []error{errors.New("totally bad thing that happened")},
			wantStopErrs:  []error{errors.New("why doesn't anything work")},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r := rebar.New(tc.givenOptions)
			r.AddProcessor(&tc.mockProcessor)
			errs := r.StartProcessors()
			assert.ElementsMatch(t, tc.wantStartErrs, errs)

			var wg sync.WaitGroup
			errs = r.StopProcessors(&wg)
			wg.Wait()
			assert.ElementsMatch(t, tc.wantStopErrs, errs)
			assert.Equal(t, 1, tc.mockProcessor.starts)
			assert.Equal(t, 1, tc.mockProcessor.stops)
		})
	}
}
