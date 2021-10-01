package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/masonhubco/rebar/v2"
	"github.com/masonhubco/rebar/v2/middleware"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_Recovery(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		shouldPanic  bool
		wantRespCode int
		wantErrCount int
	}{
		{
			name:         "should panic",
			shouldPanic:  true,
			wantRespCode: http.StatusInternalServerError,
			wantErrCount: 1,
		},
		{
			name:         "no panic",
			shouldPanic:  false,
			wantRespCode: http.StatusOK,
			wantErrCount: 0,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			resp := httptest.NewRecorder()
			_, r := gin.CreateTestContext(resp)

			r.Use(func(c *gin.Context) {
				c.Set(rebar.LoggerKey, zap.NewNop())
				c.Next()
				assert.Equal(t, tc.wantErrCount, len(c.Errors))
				assert.Equal(t, tc.wantRespCode, c.Writer.Status())
			})
			r.Use(middleware.Recovery())
			r.GET("/", func(c *gin.Context) {
				if tc.shouldPanic {
					panic("expected unit test panic")
				}
			})
			r.ServeHTTP(resp, req)
		})
	}
}
