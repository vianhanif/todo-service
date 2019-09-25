package middleware

import (
	"net/http"
	"strings"

	"github.com/vianhanif/todo-service/config"
	"github.com/vianhanif/todo-service/errors"
	"github.com/vianhanif/todo-service/internal/server/http/response"
)

// Private is a middleware for checking private api access
func Private(config *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			splitToken := strings.Split(token, " ")

			if len(splitToken) < 2 {
				response.WithError(w, nil, errors.NewAuthError("private_token_not_found"))
				return
			}

			token = strings.Trim(splitToken[1], " ")
			serverPrivateToken := config.PrivateToken

			if token != serverPrivateToken {
				response.WithError(w, nil, errors.NewAuthError("private_token_invalid"))
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
