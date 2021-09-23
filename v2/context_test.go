package rebar_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/masonhubco/rebar/v2"
	"github.com/masonhubco/rebar/v2/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
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

func Test_LoggerFrom(t *testing.T) {
	t.Parallel()

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
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := rebar.LoggerFrom(tc.given)
			assert.IsType(t, tc.want, got)
		})
	}
}

func Test_RequestIDFrom(t *testing.T) {
	t.Parallel()

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
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := rebar.RequestIDFrom(tc.given)
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_I18nFrom(t *testing.T) {
	t.Parallel()

	lang := rebar.LanguageScoped{
		Language: "es",
	}

	tests := []struct {
		name     string
		mock     func(*gin.Context)
		wantLang rebar.LanguageScoped
		isItOk   bool
	}{
		{
			name: "context does not have i18n",
			mock: func(ctx *gin.Context) {
				// doing nothing here
			},
			wantLang: rebar.LanguageScoped{},
			isItOk:   false,
		},
		{
			name: "context has i18n but it is not LanguageScoped",
			mock: func(ctx *gin.Context) {
				ctx.Set(rebar.I18nKey, "not a LanguageScoped")
			},
			wantLang: rebar.LanguageScoped{},
			isItOk:   false,
		},
		{
			name: "happy path and context has tx",
			mock: func(ctx *gin.Context) {
				ctx.Set(rebar.I18nKey, lang)
			},
			wantLang: lang,
			isItOk:   true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			resp := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(resp)
			tc.mock(ctx)

			gotLang, ok := rebar.I18nFrom(ctx)

			require.Equal(t, tc.isItOk, ok)
			assert.Equal(t, tc.wantLang, gotLang)
		})
	}
}

func Test_I18nMustFrom(t *testing.T) {
	t.Parallel()

	lang := rebar.LanguageScoped{
		Language: "es",
	}

	tests := []struct {
		name      string
		givenCtx  *gin.Context
		wantPanic bool
	}{
		{
			name: "should panic",
			givenCtx: func() *gin.Context {
				resp := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(resp)
				// not setting any value for rebar.I18nKey
				return ctx
			}(),
			wantPanic: true,
		},
		{
			name: "no panic",
			givenCtx: func() *gin.Context {
				resp := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(resp)
				ctx.Set(rebar.I18nKey, lang)
				return ctx
			}(),
			wantPanic: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if tc.wantPanic {
				assert.Panics(t, func() {
					rebar.I18nMustFrom(tc.givenCtx)
				})
			} else {
				assert.NotPanics(t, func() {
					rebar.I18nMustFrom(tc.givenCtx)
				})
			}
		})
	}
}

func Test_TxFrom(t *testing.T) {
	t.Parallel()

	tx := new(sqlx.Tx)

	tests := []struct {
		name   string
		mock   func(*gin.Context)
		wantTx *sqlx.Tx
		isItOk bool
	}{
		{
			name: "context does not have tx",
			mock: func(ctx *gin.Context) {
				// doing nothing here
			},
			wantTx: nil,
			isItOk: false,
		},
		{
			name: "context has tx but it is not sqlx tx",
			mock: func(ctx *gin.Context) {
				ctx.Set(rebar.TxKey, "not a sqlx tx")
			},
			wantTx: nil,
			isItOk: false,
		},
		{
			name: "happy path and context has tx",
			mock: func(ctx *gin.Context) {
				ctx.Set(rebar.TxKey, tx)
			},
			wantTx: tx,
			isItOk: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			resp := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(resp)
			tc.mock(ctx)

			gotTx, ok := rebar.TxFrom(ctx)

			require.Equal(t, tc.isItOk, ok)
			assert.Equal(t, tc.wantTx, gotTx)
		})
	}
}

func Test_TxMustFrom(t *testing.T) {
	t.Parallel()

	tx := new(sqlx.Tx)

	tests := []struct {
		name      string
		givenCtx  *gin.Context
		wantPanic bool
	}{
		{
			name: "should panic",
			givenCtx: func() *gin.Context {
				resp := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(resp)
				// not setting any value for rebar.TxKey
				return ctx
			}(),
			wantPanic: true,
		},
		{
			name: "no panic",
			givenCtx: func() *gin.Context {
				resp := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(resp)
				ctx.Set(rebar.TxKey, tx)
				return ctx
			}(),
			wantPanic: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if tc.wantPanic {
				assert.Panics(t, func() {
					rebar.TxMustFrom(tc.givenCtx)
				})
			} else {
				assert.NotPanics(t, func() {
					rebar.TxMustFrom(tc.givenCtx)
				})
			}
		})
	}
}
