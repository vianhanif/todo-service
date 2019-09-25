// template version: 1.0.9
package todo

import (
	"context"

	"github.com/vianhanif/todo-service/internal/storage"
)

// Service ,
type Service struct {
	q storage.Queryable
	s *Storage
}

// Single , find single Todo record matching with condition specified by query and args.
func (s *Service) Single(ctx context.Context, query string, args ...interface{}) (*Todo, error) {
	return s.s.Single(ctx, query, args...)
}

// First , find first Todo record matching with condition specified by query and args.
func (s *Service) First(ctx context.Context, query string, args ...interface{}) (*Todo, error) {
	return s.s.First(ctx, query, args...)
}

// FirstOrder , find first Todo record matching with condition specified by query and args, and ordered.
func (s *Service) FirstOrder(ctx context.Context, query, order string, args ...interface{}) (*Todo, error) {
	return s.s.FirstOrder(ctx, query, order, args...)
}

// Where , find all Todo records matching with condition specified by query and args.
func (s *Service) Where(ctx context.Context, query string, args ...interface{}) ([]*Todo, error) {
	return s.s.Where(ctx, query, args...)
}

// WhereOrder , find all Todo records matching with condition specified by query and args, and ordered.
func (s *Service) WhereOrder(ctx context.Context, query, order string, args ...interface{}) ([]*Todo, error) {
	return s.s.WhereOrder(ctx, query, order, args...)
}

// WhereWithPaging , find all Todo records matching with condition specified by query and args limiting the result specified by size
// when size has value less than 1, the function will use default value 20 for size.
// when page has value less than 1, the function will use default value 1 for page. page has base index 1
func (s *Service) WhereWithPaging(ctx context.Context, page, size int, query, order string, args ...interface{}) ([]*Todo, error) {
	return s.s.WhereWithPaging(ctx, page, size, query, order, args...)
}

// WhereNoFilter , find all Todo records matching with condition specified by query and args.
func (s *Service) WhereNoFilter(ctx context.Context, query string, args ...interface{}) ([]*Todo, error) {
	return s.s.WhereNoFilter(ctx, query, args...)
}

// FindAll , find all Todo records.
func (s *Service) FindAll(ctx context.Context, page, size int, order string) ([]*Todo, error) {
	return s.s.FindAll(ctx, page, size, order)
}

// FindByKeys , find Todo using it's primary key(s).
func (s *Service) FindByKeys(ctx context.Context, id int) (*Todo, error) {
	return s.s.FindByKeys(ctx, id)
}

// Create , create new Todo record.
func (s *Service) Create(ctx context.Context, pTodo *Todo) error {
	return s.s.Create(ctx, pTodo)
}

// Update , update Todo record.
func (s *Service) Update(ctx context.Context, pTodo *Todo) error {
	return s.s.Update(ctx, pTodo)
}

// Save , create new Todo if it doesn't exist or update if exists.
func (s *Service) Save(ctx context.Context, pTodo *Todo) error {
	return s.s.Save(ctx, pTodo)
}

// FindByKeysNoFilter , find Todo using it's primary key(s).
func (s *Service) FindByKeysNoFilter(ctx context.Context, id int) (*Todo, error) {
	return s.s.FindByKeysNoFilter(ctx, id)
}

// NewService , returns new Service.
func NewService(q storage.Queryable) *Service {
	s := NewStorage(q)
	service := &Service{
		q: q,
		s: s,
	}
	return service
}

type key int

const ctxKey key = 0

// NewContext , return new context with s.
func NewContext(ctx context.Context, s *Service) context.Context {
	return context.WithValue(ctx, ctxKey, s)
}

// FromContext , return a service from a context.
func FromContext(ctx context.Context) (*Service, bool) {
	service, ok := ctx.Value(ctxKey).(*Service)
	return service, ok
}
