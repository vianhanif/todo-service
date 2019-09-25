package middleware

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
)

// Recoverer is a middleware that recovers from panics & logs the panic &
// returns a HTTP 500 (Internal Server Error) status if possible.
func Recoverer() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					fmt.Fprintf(os.Stderr, "Panic: %+v\n", rvr)
					debug.PrintStack()
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
