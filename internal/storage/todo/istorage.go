// template version: 1.0.9
package todo

import (
	"context"
)

// IStorage interface that wraps methods for working with table todo
type IStorage interface {
	// Single , find single Todo record matching with condition specified by query and args.
	Single(ctx context.Context, query string, args ...interface{}) (*Todo, error)
	// First , find first Todo record matching with condition specified by query and args.
	First(ctx context.Context, query string, args ...interface{}) (*Todo, error)
	// FirstOrder , find first Todo record matching with condition specified by query and args, and ordered.
	FirstOrder(ctx context.Context, query, order string, args ...interface{}) (*Todo, error)
	// Where , find all Todo records matching with condition specified by query and args.
	Where(ctx context.Context, query string, args ...interface{}) ([]*Todo, error)
	// WhereOrder , find all Todo records matching with condition specified by query and args, and ordered.
	WhereOrder(ctx context.Context, query, order string, args ...interface{}) ([]*Todo, error)
	// WhereWithPaging , find all Todo records matching with condition specified by query and args limiting the result specified by size
	// when size has value less than 1, the function will use default value 20 for size.
	// when page has value less than 1, the function will use default value 1 for page. page has base index 1
	WhereWithPaging(ctx context.Context, page, size int, query, order string, args ...interface{}) ([]*Todo, error)
	// WhereNoFilter , find all Todo records matching with condition specified by query and args.
	WhereNoFilter(ctx context.Context, query string, args ...interface{}) ([]*Todo, error)
	// FindAll , find all Todo records.
	FindAll(ctx context.Context) ([]*Todo, error)
	// FindByKeys , find Todo using it's primary key(s).
	FindByKeys(ctx context.Context, id int) (*Todo, error)
	// FindByKeysNoFilter , find Todo using it's primary key(s).
	FindByKeysNoFilter(ctx context.Context, id int) (*Todo, error)
	// Create , create new Todo record.
	Create(ctx context.Context, pTodo *Todo) error
	// Update , update Todo record.
	Update(ctx context.Context, pTodo *Todo) error
	// Delete , remove Todo using it's primary key(s).
	Delete(ctx context.Context, id int) error
	// Save , create new Todo if it doesn't exist or update if exists.
	Save(ctx context.Context, pTodo *Todo) error
}
