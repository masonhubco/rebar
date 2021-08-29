package rebar_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/masonhubco/rebar/v2"
	"github.com/stretchr/testify/assert"
)

func Test_Options_Mode(t *testing.T) {
	tests := []struct {
		name  string
		given string
		want  string
	}{
		{name: "development", given: rebar.Development, want: gin.DebugMode},
		{name: "test", given: rebar.Test, want: gin.TestMode},
		{name: "staging", given: rebar.Staging, want: gin.ReleaseMode},
		{name: "integration", given: rebar.Staging, want: gin.ReleaseMode},
		{name: "production", given: rebar.Staging, want: gin.ReleaseMode},
		{name: "unknown", given: "unknown", want: gin.ReleaseMode},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			opts := rebar.Options{
				Environment: tc.given,
			}
			assert.Equal(t, tc.want, opts.Mode())
		})
	}
}
