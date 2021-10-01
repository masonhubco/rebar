package rebar_test

import (
	"bytes"
	"io"
	"os"
	"sync"
	"testing"

	"github.com/masonhubco/rebar/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewStandardLogger(t *testing.T) {
	tests := []struct {
		name string
		log  func(rebar.Logger)
		want []string
	}{
		{
			name: "log info",
			log: func(lg rebar.Logger) {
				lg.Info("some info")
			},
			want: []string{"INFO:", "some info"},
		},
		{
			name: "log warning",
			log: func(lg rebar.Logger) {
				lg.Warn("some warning")
			},
			want: []string{"WARNING:", "some warning"},
		},
		{
			name: "log error",
			log: func(lg rebar.Logger) {
				lg.Error("some error")
			},
			want: []string{"ERROR:", "some error"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			captured := stdoutAndStderrFrom(t, func() {
				logger, err := rebar.NewStandardLogger()
				require.NoError(t, err)
				tc.log(logger)
			})
			for _, want := range tc.want {
				assert.Contains(t, captured, want)
			}
		})
	}
}

func stdoutAndStderrFrom(t *testing.T, fn func()) string {
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
	}()
	os.Stdout = writer
	os.Stderr = writer
	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		io.Copy(&buf, reader)
		out <- buf.String()
	}()
	wg.Wait()
	fn()
	writer.Close()
	return <-out
}
