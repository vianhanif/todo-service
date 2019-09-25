package v1_test

import (
	"context"
	"testing"

	"github.com/vianhanif/go-pkg/sql/helper"
	"github.com/vianhanif/todo-service/internal/app"
	"github.com/vianhanif/todo-service/internal/test"
)

var (
	ctx context.Context
	a   *app.App
)

func init() {
	ctx = context.TODO()
	a = test.GetApp(ctx)
}

func TestCreate(t *testing.T) {
	item, err := a.Services.Todo.Create(ctx, test.FakeTodo())
	if err != nil {
		t.Fatal(err)
	}
	if item == nil {
		t.Fatal("data is nil")
	}
}

func TestList(t *testing.T) {
	items, err := a.Services.Todo.List(ctx, []helper.QueryFilter{
		helper.QueryFilter{Key: "order", Operation: "order", Column: `"created_at"`, Value: `DESC`},
		helper.QueryFilter{Key: "page", Operation: "offset", Value: "1"},
		helper.QueryFilter{Key: "perPage", Operation: "limit", Value: "10"},
	}...)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) == 0 {
		t.Fatal("data is empty")
	}
}

func TestUpdate(t *testing.T) {
	newTitle := "new title"
	item, _ := a.Services.Todo.Create(ctx, test.FakeTodo())
	item.Title = newTitle
	item, err := a.Services.Todo.Update(ctx, item)
	if err != nil {
		t.Fatal(err)
	}
	if item == nil {
		t.Fatal("data is nil")
	}
	if item.Title != newTitle {
		t.Fatal("title not changed")
	}
}

func TestDelete(t *testing.T) {
	item, _ := a.Services.Todo.Create(ctx, test.FakeTodo())
	err := a.Services.Todo.Delete(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
}
