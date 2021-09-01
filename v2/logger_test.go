package rebar_test

import (
	"bytes"
	"io"
	"os"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/masonhubco/rebar/v2"
	"github.com/masonhubco/rebar/v2/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
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

func Test_LoggerFrom(t *testing.T) {
	tests := []struct {
		name  string
		given *gin.Context
		want  rebar.Logger
	}{
		{
			name: "no logger in context",
			given: &gin.Context{
				Keys: map[string]interface{}{
					"notALogger": "not a logger",
				},
			},
			want: &zap.Logger{}, // only need to check the type
		},
		{
			name: "logger in context is not an instance of Logger",
			given: &gin.Context{
				Keys: map[string]interface{}{
					rebar.LoggerKey: "not an instance of Logger",
				},
			},
			want: &zap.Logger{}, // only need to check the type
		},
		{
			name: "happy path",
			given: &gin.Context{
				Keys: map[string]interface{}{
					rebar.LoggerKey: &mocks.Logger{},
				},
			},
			want: &mocks.Logger{}, // only need to check the type
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := rebar.LoggerFrom(tc.given)
			assert.IsType(t, tc.want, got)
		})
	}
}

func Test_RequestIDFrom(t *testing.T) {
	tests := []struct {
		name  string
		given *gin.Context
		want  string
	}{
		{
			name: "no request id in context",
			given: &gin.Context{
				Keys: map[string]interface{}{
					"notARequestID": "not a request ID",
				},
			},
			want: "",
		},
		{
			name: "request id in context is not a string",
			given: &gin.Context{
				Keys: map[string]interface{}{
					rebar.RequestIDKey: 12345, // wrong type should be a string
				},
			},
			want: "",
		},
		{
			name: "happy path",
			given: &gin.Context{
				Keys: map[string]interface{}{
					rebar.RequestIDKey: "test-request-id",
				},
			},
			want: "test-request-id",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := rebar.RequestIDFrom(tc.given)
			assert.Equal(t, tc.want, got)
		})
	}
}
