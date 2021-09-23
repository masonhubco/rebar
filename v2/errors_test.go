package rebar_test

import (
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/masonhubco/rebar/v2"
	"github.com/stretchr/testify/assert"
)

func Test_ContextErrors_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		given []*gin.Error
		want  string
	}{
		{
			name:  "empty error",
			given: []*gin.Error{},
			want:  "",
		},
		{
			name: "errors without meta",
			given: []*gin.Error{
				{
					Err:  errors.New("expected unit test error"),
					Type: gin.ErrorTypePrivate,
				},
			},
			want: "Error #01: expected unit test error\n",
		},
		{
			name: "errors with meta",
			given: []*gin.Error{
				{
					Err:  errors.New("expected unit test error"),
					Type: gin.ErrorTypePrivate,
					Meta: gin.H{"some": "meta data"},
				},
			},
			want: "Error #01: expected unit test error\n     Meta: map[some:meta data]\n",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			errs := rebar.NewContextErrors(tc.given)
			assert.Equal(t, tc.want, errs.Error())
		})
	}
}

func Test_ContextErrors_Is(t *testing.T) {
	t.Parallel()

	err := errors.New("expected unit test error")
	ginErr := &gin.Error{
		Err:  err,
		Type: gin.ErrorTypePrivate,
	}
	ctxErrs := rebar.NewContextErrors([]*gin.Error{ginErr})

	tests := []struct {
		name         string
		givenCtxErrs *rebar.ContextErrors
		givenTarget  error
		want         bool
	}{
		{
			name:         "yes and compare a target to itself",
			givenCtxErrs: ctxErrs,
			givenTarget:  ctxErrs,
			want:         true,
		},
		{
			name:         "yes and target is a gin error",
			givenCtxErrs: ctxErrs,
			givenTarget:  ginErr,
			want:         true,
		},
		{
			name:         "yes and target is not a gin error",
			givenCtxErrs: ctxErrs,
			givenTarget:  err,
			want:         true,
		},
		{
			name:         "no because target is not a part of the context errors",
			givenCtxErrs: ctxErrs,
			givenTarget:  errors.New("a different error"),
			want:         false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.want, tc.givenCtxErrs.Is(tc.givenTarget))
		})
	}
}
