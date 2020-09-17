package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/masonhubco/rebar/middleware"
	"github.com/stretchr/testify/assert"
)

func dummyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func Test_AuthenticationMW(t *testing.T) {
	router := mux.NewRouter()

	auth := middleware.AuthenticationMW{SystemToken: "token"}
	router.Use(auth.Authenticate)
	router.HandleFunc("/", dummyHandler).Methods("GET")

	t.Run("bad token", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", rr.Body)
		req.Header.Add("Authorization", "Bearer bad token")
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("no bearer", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", rr.Body)
		req.Header.Add("Authorization", "bad token")
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("no token", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", rr.Body)
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("good token", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", rr.Body)
		req.Header.Add("Authorization", "Bearer token")
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

}
