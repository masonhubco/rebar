package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/masonhubco/rebar/v2"
	"github.com/masonhubco/rebar/v2/middleware"
	"github.com/stretchr/testify/assert"
)

func Test_ForceSSL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		givenEnv    string
		wantAborted bool
	}{
		{
			name:        "no redirect",
			givenEnv:    rebar.Development,
			wantAborted: false,
		},
		{
			name:        "redirect",
			givenEnv:    rebar.Production,
			wantAborted: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			resp := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(resp)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)

			middleware.ForceSSL(tc.givenEnv)(ctx)

			assert.Equal(t, tc.wantAborted, ctx.IsAborted())
		})
	}
}
