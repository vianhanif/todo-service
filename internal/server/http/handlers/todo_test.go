package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-chi/chi"

	"github.com/vianhanif/todo-service/config"
	"github.com/vianhanif/todo-service/internal/app"
	"github.com/vianhanif/todo-service/internal/server/http/handlers"
	"github.com/vianhanif/todo-service/internal/server/http/middleware"
	"github.com/vianhanif/todo-service/internal/storage/todo"
	"github.com/vianhanif/todo-service/internal/test"
)

var a *app.App
var ctx context.Context

func init() {
	ctx = context.TODO()
	a = test.GetApp(ctx)
}

func TestCreate(t *testing.T) {
	item := test.FakeTodo()
	jsonBody := map[string]interface{}{
		"title":  item.Title,
		"detail": item.Detail,
	}
	body, _ := json.Marshal(jsonBody)

	handler := handlerFunc(middleware.Private(a.Config), handlers.Create())
	req := private(a.Config)(httpRequestor(
		"/private/v1/todos",
		"POST",
		handler,
		body,
	))

	res := httpWriter()
	handler.ServeHTTP(res, req)

	status := res.Code
	if status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	data := &todo.Todo{}
	err := json.NewDecoder(res.Body).Decode(data)
	if err != nil {
		t.Error("error decoding body response")
	}
	if data.Title != item.Title {
		t.Error("title not correct")
	}
}

func TestList(t *testing.T) {
	handler := handlerFunc(middleware.Private(a.Config), handlers.List())
	req := private(a.Config)(httpRequestor(
		"/private/v1/todos",
		"GET",
		handler,
		nil,
	))

	res := httpWriter()
	handler.ServeHTTP(res, req)

	status := res.Code
	if status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	data := []*todo.Todo{}
	err := json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		t.Error("error decoding body response")
	}
	if len(data) == 0 {
		t.Error("data empty")
	}
}

func TestFind(t *testing.T) {
	item, err := a.Services.Todo.Create(ctx, test.FakeTodo())
	if err != nil {
		t.Fatal(err)
	}

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", strconv.Itoa(item.ID))
	ctx := context.WithValue(context.TODO(), chi.RouteCtxKey, rctx)

	handler := handlerFunc(middleware.Private(a.Config), handlers.Find())
	req := private(a.Config)(httpRequestor(
		fmt.Sprintf("/private/v1/todos/%d", item.ID),
		"GET",
		handler,
		nil,
	))

	res := httpWriter()
	handler.ServeHTTP(res, req.WithContext(ctx))

	status := res.Code
	if status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	data := &todo.Todo{}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		t.Error("error decoding body response")
	}
	if data.ID != item.ID {
		t.Error("error get data")
	}
}

func httpRequestor(path, method string, handler http.Handler, body []byte) *http.Request {
	return httptest.NewRequest(method, path, bytes.NewReader(body))
}

func httpWriter() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}

func handlerFunc(middleware func(next http.Handler) http.Handler, h http.HandlerFunc) http.Handler {
	return app.InjectorMiddleware(a)(middleware(http.HandlerFunc(h)))
}

func private(config *config.Config) func(req *http.Request) *http.Request {
	return func(req *http.Request) *http.Request {
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf(`Bearer %s`, config.PrivateToken))
		return req
	}
}
