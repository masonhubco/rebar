package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/masonhubco/rebar/v2"
	"github.com/masonhubco/rebar/v2/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_I18n(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		givenLang string
		wantLang  rebar.LanguageScoped
	}{
		{
			name:      "accept language header not set",
			givenLang: "",
			wantLang: rebar.LanguageScoped{
				Language: "en",
			},
		},
		{
			name:      "accept language header set to fr",
			givenLang: "fr",
			wantLang: rebar.LanguageScoped{
				Language: "fr",
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			resp := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(resp)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
			if tc.givenLang != "" {
				ctx.Request.Header.Set("Accept-Language", tc.givenLang)
			}

			middleware.I18n()(ctx)
			gotLang, ok := rebar.I18nFrom(ctx)

			require.True(t, ok)
			assert.Equal(t, tc.wantLang, gotLang)
		})
	}
}
