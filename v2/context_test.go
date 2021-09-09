package rebar_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/masonhubco/rebar/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testBuffaloValidateError struct {
	Errors                     map[string][]string `json:"errors"`
	rebar.BuffaloValidateError `json:"-"`
}

func newTestBuffaloVlidateError() error {
	return &testBuffaloValidateError{
		Errors: map[string][]string{
			"input_field": {"cannot be left blank"},
		},
	}
}

func Test_AbortWithError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		givenCode int
		givenErr  error
		wantBody  string
	}{
		{
			name:      "with buffalo validate error",
			givenCode: http.StatusBadRequest,
			givenErr:  newTestBuffaloVlidateError(),
			wantBody:  `{"request_id":"test-request-id","validate":{"errors":{"input_field":["cannot be left blank"]}}}`,
		},
		{
			name:      "with just an error",
			givenCode: http.StatusTeapot,
			givenErr:  errors.New("expected unit test error"),
			wantBody:  `{"request_id":"test-request-id","error":"expected unit test error"}`,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			resp := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(resp)
			ctx.Set(rebar.RequestIDKey, "test-request-id")

			rebar.AbortWithError(ctx, tc.givenCode, tc.givenErr)

			require.Equal(t, tc.givenCode, resp.Code)
			assert.JSONEq(t, tc.wantBody, resp.Body.String())
		})
	}
}
