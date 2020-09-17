package middleware

import (
	"log"
	"net/http"
)

//Logger is middleware that wraps the request and logs out
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request uri: %s url: %s", r.RequestURI, r.URL)
		next.ServeHTTP(w, r)
	})
}
