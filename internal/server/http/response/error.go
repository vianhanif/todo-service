package response

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/vianhanif/todo-service/errors"
	"github.com/vianhanif/todo-service/internal/alert"
)

//ErrorResponse represents error message
//swagger:model
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// WithError emits proper error response
func WithError(w http.ResponseWriter, a alert.Alert, err error) {
	switch err.(type) {
	case errors.BaseError:
		WithData(w, a, http.StatusBadRequest, err)
	case errors.NotFoundError:
		WithData(w, a, http.StatusNotFound, err)
	case errors.CommonError:
		WithData(w, a, http.StatusBadRequest, err)
	case errors.ValidationError:
		WithData(w, a, http.StatusUnprocessableEntity, err)
	case errors.AuthError:
		WithData(w, a, http.StatusUnauthorized, err)
	case errors.ServiceError:
		logError(a, err)
		response := ErrorResponse{
			Code:    http.StatusText(http.StatusInternalServerError),
			Message: "Server tidak dapat memproses permintaan anda, cobalah beberapa saat lagi.",
		}
		WithData(w, a, http.StatusInternalServerError, response)
	default:
		logError(a, err)
		response := ErrorResponse{
			Code:    http.StatusText(http.StatusInternalServerError),
			Message: "Server tidak dapat memproses permintaan anda, cobalah beberapa saat lagi.",
		}
		WithData(w, a, http.StatusInternalServerError, response)
	}
}

// WithData emits response with data
func WithData(w http.ResponseWriter, a alert.Alert, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func logError(a alert.Alert, err error) {
	msg := fmt.Sprintf("%+v\n%s", err, string(debug.Stack()))
	log.Println(msg)
	if a != nil {
		alert := alert.NewAlert(err.Error(), err, debug.Stack())
		a.Alert(alert)
	}
}
