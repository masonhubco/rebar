package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

//AuthenticationMW allows you to instantiate with the SystemToken and pass it into the
//middleware function, thereby injecting your runtime dependency
type AuthenticationMW struct {
	SystemToken string
}

// Authenticate will ensure that the user has permissions to access the respective handler
func (amw AuthenticationMW) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := amw.getCallbackToken(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}

		if tokenString != amw.SystemToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

// getCallbackToken gets the token from the Authorization header
// removes the Bearer part from the authorisation header value.
// returns No token error if Token is not found
// returns Token Invalid error if the token value cannot be obtained by removing `Bearer `
func (amw AuthenticationMW) getCallbackToken(authString string) (string, error) {
	if authString == "" {
		return "", fmt.Errorf("No token provided")
	}
	splitToken := strings.Split(authString, "Bearer ")
	if len(splitToken) != 2 {
		return "", fmt.Errorf("what kind of token is this...it's terribly formatted")
	}
	tokenString := splitToken[1]
	return tokenString, nil
}
