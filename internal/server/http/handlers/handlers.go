package handlers

import (
	"net/http"

	"github.com/vianhanif/todo-service/internal/alert"
	"github.com/vianhanif/todo-service/internal/app"
	"github.com/vianhanif/todo-service/internal/server/http/response"
)

// HandlerFunc .
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// Handler .
func Handler(fn HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			return
		default:
			err := fn(w, r)
			if err != nil {
				app := app.FromContext(r.Context())
				var alert alert.Alert
				if app != nil {
					alert = app.Alert
				}
				response.WithError(w, alert, err)
			}
			return
		}
	}
}
