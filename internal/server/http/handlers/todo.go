package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/vianhanif/todo-service/errors"
	"github.com/vianhanif/todo-service/internal/app"
	"github.com/vianhanif/todo-service/internal/server/http/response"
	"github.com/vianhanif/todo-service/internal/storage/todo"
)

// Create .
func Create() http.HandlerFunc {
	type payload struct {
		Title  string `json:"title"`
		Detail string `json:"detail"`
	}

	return Handler(func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		app := app.FromContext(ctx)

		var params payload

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			return errors.NewCommonError("Terjadi kesalahan data. Silahkan mencoba kembali")
		}

		data := &todo.Todo{
			Title:  params.Title,
			Detail: params.Detail,
		}
		record, err := app.Services.Todo.Create(ctx, data)
		if err != nil {
			return errors.NewServiceError("Proses gagal. Silahkan mencoba beberapa saat lagi")
		}

		response.JSON(w, http.StatusOK, record)
		return nil
	})
}
