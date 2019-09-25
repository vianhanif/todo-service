package v1

import (
	"context"

	"github.com/vianhanif/go-pkg/sql/helper"

	"github.com/vianhanif/todo-service/internal/storage"
	"github.com/vianhanif/todo-service/internal/storage/todo"
)

// Service .
type Service interface {
	Create(ctx context.Context, item *todo.Todo) (record *todo.Todo, err error)
	List(ctx context.Context, filters ...helper.QueryFilter) (records []*todo.Todo, err error)
	Find(ctx context.Context, filters ...helper.QueryFilter) (records *todo.Todo, err error)
	Update(ctx context.Context, item *todo.Todo) (record *todo.Todo, err error)
	Delete(ctx context.Context, ID int) error
}

// App .
type App struct {
	query storage.Queryable
	todo  todo.IStorage
}

// Create .
func (s *App) Create(ctx context.Context, item *todo.Todo) (record *todo.Todo, err error) {
	return item, s.todo.Create(ctx, item)
}

// List .
func (s *App) List(ctx context.Context, filters ...helper.QueryFilter) ([]*todo.Todo, error) {
	where, _ := helper.BuildFilter(filters...)
	return s.todo.FindAll(ctx, where.Offset(), where.Limit(), where.OrderBy())
}

// Find .
func (s *App) Find(ctx context.Context, filters ...helper.QueryFilter) (record *todo.Todo, err error) {
	where, args := helper.BuildFilter(filters...)
	return s.todo.Single(ctx, where.String(), args...)
}

// Update .
func (s *App) Update(ctx context.Context, item *todo.Todo) (*todo.Todo, error) {
	return item, s.todo.Update(ctx, item)
}

// Delete .
func (s *App) Delete(ctx context.Context, ID int) error {
	return s.todo.Delete(ctx, ID)
}

// NewService .
func NewService(query storage.Queryable, t *todo.Storage) Service {
	return &App{
		query: query,
		todo:  t,
	}
}
