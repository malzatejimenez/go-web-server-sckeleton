package middlewares

import (
	"net/http"
	"platzi/go/rest-ws/models"
	"platzi/go/rest-ws/server"
	"strings"

	"github.com/golang-jwt/jwt"
)

// NO_AUTH_NEEDED is a list of paths that don't need authentication
var (
	NO_AUTH_NEEDED = []string{"/", "/signup", "/login"}
)

// shouldCheckToken is a function that checks if the token should be checked
func shouldCheckToken(route string) bool {
	// check if the route is in the list of routes that don't need authentication
	for _, path := range NO_AUTH_NEEDED {
		// if the route is in the list, return false
		if path == route {
			return false
		}
	}

	// return true if the route is not in the list
	return true
}

// CheckAuthMiddleware is a middleware that checks if the user is authenticated
func CheckAuthMiddleware(s server.Server) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// check if the token should be checked
			if !shouldCheckToken(r.URL.Path) {
				// if the token should not be checked, call the next handler
				next.ServeHTTP(w, r)

				// don't need to continue with the execution
				return
			}

			// get the token from the request
			tokenString := strings.TrimSpace(r.Header.Get("Authorization"))

			// check if the token is empty
			if tokenString == "" {
				// if the token is empty, return an error
				http.Error(w, "Unauthorized", http.StatusUnauthorized)

				// don't need to continue with the execution
				return
			}

			// validate the token
			_, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(t *jwt.Token) (interface{}, error) {
				return []byte(s.Config().JwtSecret), nil
			})

			// check if there was an error
			if err != nil {
				// if there was an error, return an error
				http.Error(w, "Unauthorized", http.StatusUnauthorized)

				// don't need to continue with the execution
				return
			}

			// if there was no error, call the next handler
			next.ServeHTTP(w, r)
		})
	}
}
