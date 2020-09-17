package helpers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/masonhubco/mercury/rebar"
	"github.com/stretchr/testify/assert"
)

func Test_APIRender(t *testing.T) {
	t.Run("Good Widget", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", rr.Body)
		type Widgets struct {
			FieldOne string `json:"FieldOne"`
			FieldTwo string `json:"FieldTwo"`
		}

		rebar.APIRenderList(rr, req, 1, Widgets{FieldOne: "one", FieldTwo: "two"})
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, `{"FieldOne":"one","FieldTwo":"two"}
`, rr.Body.String())
		assert.Equal(t, "1", rr.Header().Get("Result-Count"))
	})

	t.Run("Bad Widget", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", rr.Body)
		type Widgets struct {
			FieldOne func()
		}

		rebar.APIRenderList(rr, req, 1, Widgets{FieldOne: func() {}})
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, "1", rr.Header().Get("Result-Count"))
		assert.Equal(t, "error rendering JSON\n", rr.Body.String())

	})

}
