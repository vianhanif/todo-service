package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

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
	ctx := context.TODO()
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
