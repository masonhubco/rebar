package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/masonhubco/rebar/v2/middleware"
	"github.com/masonhubco/rebar/v2/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_Transaction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		givenCtx    *gin.Context
		givenDBErr  error
		wantAborted bool
		wantCtxErr  string
		wantStatus  int
	}{
		{
			name: "has database commit error",
			givenCtx: func() *gin.Context {
				resp := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(resp)
				return ctx
			}(),
			givenDBErr:  errors.New("expected database commit error"),
			wantAborted: true,
			wantCtxErr:  "Error #01: expected database commit error\n",
			wantStatus:  http.StatusInternalServerError,
		},
		{
			name: "has gin context error",
			givenCtx: func() *gin.Context {
				resp := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(resp)
				ctx.Status(http.StatusUnprocessableEntity)
				ctx.Error(errors.New("expected gin context error"))
				return ctx
			}(),
			givenDBErr:  nil, // no database commit error
			wantAborted: true,
			wantCtxErr:  "Error #01: expected gin context error\n",
			wantStatus:  http.StatusUnprocessableEntity,
		},
		{
			name: "no context error but has http error code",
			givenCtx: func() *gin.Context {
				resp := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(resp)
				ctx.JSON(http.StatusTeapot, gin.H{
					"error": "expected http error",
				})
				return ctx
			}(),
			givenDBErr:  nil,   // no database commit error
			wantAborted: false, // gin context will not be aborted
			wantCtxErr:  "",
			wantStatus:  http.StatusTeapot,
		},
		{
			name: "happy path 200 ok",
			givenCtx: func() *gin.Context {
				resp := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(resp)
				ctx.JSON(http.StatusOK, gin.H{
					"some": "data returned",
				})
				return ctx
			}(),
			givenDBErr:  nil,   // no database commit error
			wantAborted: false, // gin context will not be aborted
			wantCtxErr:  "",
			wantStatus:  http.StatusOK,
		},
	}

	type withTxFn func(*sqlx.Tx) error

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			database := mocks.NewTxWrapper(ctrl)
			database.EXPECT().
				WithTx(nil, gomock.AssignableToTypeOf((withTxFn)(nil))).
				DoAndReturn(func(tx *sqlx.Tx, fn withTxFn) error {
					// simulate database commit error returned from WithTx
					if tc.givenDBErr != nil {
						return tc.givenDBErr
					}
					// run withTxFn and return its error
					// it could be either gin context error or errNonSuccess
					return fn(new(sqlx.Tx))
				})

			middleware.Transaction(database)(tc.givenCtx)

			assert.Equal(t, tc.wantAborted, tc.givenCtx.IsAborted())
			assert.Equal(t, tc.wantCtxErr, tc.givenCtx.Errors.String())
			assert.Equal(t, tc.wantStatus, tc.givenCtx.Writer.Status())
		})
	}
}
