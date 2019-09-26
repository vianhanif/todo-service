package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/vianhanif/go-pkg/sql/helper"

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

// List .
func List() http.HandlerFunc {
	return Handler(func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		app := app.FromContext(ctx)

		page, perPage := onPagination(r.URL.Query()["page"], r.URL.Query()["perPage"])
		valOffset := strconv.Itoa(page)
		valLimit := strconv.Itoa(perPage)

		params := []helper.QueryFilter{
			helper.QueryFilter{Key: "order", Operation: "order", Column: `"created_at"`, Value: `DESC`},
			helper.QueryFilter{Key: "page", Operation: "offset", Value: valOffset},
			helper.QueryFilter{Key: "perPage", Operation: "limit", Value: valLimit},
		}

		filters, anyBadRequest := helper.GetQueries(r,
			params,     // available filters
			[]string{}, // required params
		)
		if anyBadRequest != nil {
			return anyBadRequest
		}

		records, err := app.Services.Todo.List(ctx, filters...)
		if err != nil {
			return errors.NewServiceError("Proses gagal. Silahkan mencoba beberapa saat lagi")
		}

		list := []*todo.Todo{}
		for _, record := range records {
			list = append(list, record)
		}

		response.JSON(w, http.StatusOK, list)
		return nil
	})
}

// Find .
func Find() http.HandlerFunc {
	return Handler(func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		app := app.FromContext(ctx)

		ID := chi.URLParam(r, "id")
		params := []helper.QueryFilter{
			helper.QueryFilter{Key: "id", Value: ID},
		}

		filters, anyBadRequest := helper.GetQueries(r,
			params,     // available filters
			[]string{}, // required params
		)
		if anyBadRequest != nil {
			return anyBadRequest
		}

		record, err := app.Services.Todo.Find(ctx, filters...)
		if err != nil {
			return errors.NewServiceError("Proses gagal. Silahkan mencoba beberapa saat lagi")
		}

		response.JSON(w, http.StatusOK, record)
		return nil
	})
}

func onPagination(pageQuery []string, perPageQuery []string) (int, int) {
	var page, limit int
	if len(pageQuery) > 0 {
		valPage, _ := strconv.Atoi(pageQuery[0])
		page = valPage
	} else {
		page = 1
	}
	if len(perPageQuery) > 0 {
		valPerPage, _ := strconv.Atoi(perPageQuery[0])
		limit = valPerPage
	} else {
		limit = 10
	}
	offset := (page - 1) * limit
	return offset, limit
}
