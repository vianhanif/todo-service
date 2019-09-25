// template version: 1.0.9

package todo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/vianhanif/todo-service/internal/storage"
)

// Storage ,
type Storage struct {
	q storage.Queryable
}

//rowScanner represent rows object from sql
type rowScanner interface {
	Scan(dest ...interface{}) error
}

//rowScanner represent rows object from sql
type rowsResult interface {
	Next() bool
	Scan(dest ...interface{}) error
}

// NewStorage , Create new Storage.
func NewStorage(q storage.Queryable) *Storage {
	return &Storage{
		q: q,
	}
}

// Single , find one Todo record matching with condition specified by query and args.
func (s *Storage) Single(ctx context.Context, query string, args ...interface{}) (*Todo, error) {
	q := s.pickQueryable(ctx)
	stmt := fmt.Sprintf(`%s WHERE %s LIMIT 2`, selectQuery(), query)
	rows, err := q.QueryContext(ctx, stmt, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	count := 0
	var result *Todo
	for rows.Next() {
		if count > 1 {
			return nil, errors.New("found more than one record")
		}
		data := &Todo{}
		err := scan(rows, data)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
		result = data
		count++
	}
	return result, nil
}

// First , find first Todo record matching with condition specified by query and args.
func (s *Storage) First(ctx context.Context, query string, args ...interface{}) (*Todo, error) {
	q := s.pickQueryable(ctx)
	stmt := fmt.Sprintf(`%s WHERE %s LIMIT 1`, selectQuery(), query)
	row := q.QueryRowContext(ctx, stmt, args...)

	data := &Todo{}
	err := scan(row, data)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return data, nil
}

// FirstOrder , find first Todo record matching with condition specified by query and args, and ordered.
func (s *Storage) FirstOrder(ctx context.Context, query, order string, args ...interface{}) (*Todo, error) {
	q := s.pickQueryable(ctx)
	stmt := fmt.Sprintf(`%s WHERE %s ORDER BY %s LIMIT 1`, selectQuery(), query, order)
	row := q.QueryRowContext(ctx, stmt, args...)

	data := &Todo{}
	err := scan(row, data)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return data, nil
}

// Where , find all Todo records matching with condition specified by query and args.
func (s *Storage) Where(ctx context.Context, query string, args ...interface{}) ([]*Todo, error) {
	q := s.pickQueryable(ctx)
	stmt := fmt.Sprintf(`%s WHERE (%s) AND %s`, selectQuery(), query, defaultFilter())
	rows, err := q.QueryContext(ctx, stmt, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRows(rows)
}

// WhereOrder , find all Todo records matching with condition specified by query and args.
func (s *Storage) WhereOrder(ctx context.Context, query, order string, args ...interface{}) ([]*Todo, error) {
	q := s.pickQueryable(ctx)
	stmt := fmt.Sprintf(`%s WHERE (%s) AND %s ORDER BY %s`, selectQuery(), query, defaultFilter(), order)
	rows, err := q.QueryContext(ctx, stmt, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRows(rows)
}

// WhereWithPaging , find all Todo records matching with condition specified by query and args limiting the result specified by size
// when size has value less than 1, the function will use default value 20 for size.
// when page has value less than 1, the function will use default value 1 for page. page has base index 1
func (s *Storage) WhereWithPaging(ctx context.Context, page, size int, query, order string, args ...interface{}) ([]*Todo, error) {
	q := s.pickQueryable(ctx)
	limit := size
	if limit < 1 {
		limit = 20
	}
	offset := page
	if offset < 1 {
		offset = 1
	}
	offset = (offset - 1) * limit
	stmt := fmt.Sprintf(`%s WHERE (%s) AND %s ORDER BY %s LIMIT %v OFFSET %v`, selectQuery(), query, defaultFilter(), order, limit, offset)
	rows, err := q.QueryContext(ctx, stmt, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRows(rows)
}

// WhereNoFilter , find all Todo records matching with condition specified by query and args.
func (s *Storage) WhereNoFilter(ctx context.Context, query string, args ...interface{}) ([]*Todo, error) {
	q := s.pickQueryable(ctx)
	stmt := fmt.Sprintf(`%s WHERE %s`, selectQuery(), query)
	rows, err := q.QueryContext(ctx, stmt, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRows(rows)
}

// FindAll , find all Todo records.
func (s *Storage) FindAll(ctx context.Context, page, size int, order string) ([]*Todo, error) {
	q := s.pickQueryable(ctx)
	limit := size
	if limit < 1 {
		limit = 20
	}
	offset := page
	if offset < 1 {
		offset = 1
	}
	offset = (offset - 1) * limit
	stmt := fmt.Sprintf(`%s WHERE %s ORDER BY %s LIMIT %v OFFSET %v`, selectQuery(), defaultFilter(), order, limit, offset)
	rows, err := q.QueryContext(ctx, stmt)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRows(rows)
}

// FindByKeys , find Todo using it's primary key(s).
func (s *Storage) FindByKeys(ctx context.Context, id int) (*Todo, error) {
	criteria := `"id" = $1`
	stmt := fmt.Sprintf(`(%s) AND %s`, criteria, defaultFilter())
	return s.Single(ctx, stmt, id)
}

// FindByKeysNoFilter , find Todo using it's primary key(s) without filter.
func (s *Storage) FindByKeysNoFilter(ctx context.Context, id int) (*Todo, error) {
	criteria := `"id" = $1`
	return s.Single(ctx, criteria, id)
}

// Create , create new Todo record.
func (s *Storage) Create(ctx context.Context, p *Todo) error {
	q := s.pickQueryable(ctx)
	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now

	stmt, args := InsertQuery(p)
	row := q.QueryRowContext(ctx, stmt, args...)
	return scan(row, p)
}

// Update , update Todo record.
func (s *Storage) Update(ctx context.Context, p *Todo) error {
	q := s.pickQueryable(ctx)
	record, err := s.FindByKeys(ctx,
		p.ID,
	)
	if err != nil {
		return err
	}
	if record == nil {
		return errors.New("record not found")
	}
	record.Title = p.Title
	record.Detail = p.Detail
	record.DeletedAt = p.DeletedAt
	now := time.Now().UTC()
	record.UpdatedAt = now

	stmt, args := UpdateQuery(record)
	row := q.QueryRowContext(ctx, stmt, args...)
	return scan(row, p)
}

// Delete , delete Todo using it's primary key(s).
func (s *Storage) Delete(ctx context.Context, id int) error {
	q := s.pickQueryable(ctx)
	stmt, args := deleteQuery(id)
	result, err := q.ExecContext(ctx, stmt, args...)

	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	} else if count < 1 {
		return errors.New("Delete operation expecting affected rows greater than 0")
	}
	return nil
}

// Save , create new Todo if it doesn't exist or update if exists.
func (s *Storage) Save(ctx context.Context, p *Todo) error {
	record, err := s.FindByKeys(ctx,
		p.ID,
	)
	if err != nil {
		return err
	}
	if record != nil {
		return s.Update(ctx, p)
	}
	return s.Create(ctx, p)
}

func (s *Storage) pickQueryable(ctx context.Context) storage.Queryable {
	q, ok := storage.QueryableFromContext(ctx)
	if !ok {
		q = s.q
	}
	return q
}
func fields() string {
	return `"id", "title", "detail", "created_at", "updated_at", "deleted_at"`
}

func selectQuery() string {
	return fmt.Sprintf(`SELECT %s FROM "todo"`, fields())
}

// InsertQuery returns query statement and slice of arguments
func InsertQuery(data *Todo) (string, []interface{}) {
	o := []string{
		"title",
		"detail",
		"created_at",
		"updated_at",
		"deleted_at",
	}
	m := map[string]interface{}{
		"title":      data.Title,
		"detail":     data.Detail,
		"created_at": data.CreatedAt,
		"updated_at": data.UpdatedAt,
		"deleted_at": data.DeletedAt,
	}

	fs, ph := func(v map[string]interface{}) ([]string, []string) {
		fs := []string{}
		ph := []string{}
		for i, k := range o {
			fs = append(fs, fmt.Sprintf(`"%s"`, k))
			ph = append(ph, fmt.Sprintf(`$%d`, i+1))
		}
		return fs, ph
	}(m)
	args := func(v map[string]interface{}) []interface{} {
		args := []interface{}{}
		for _, k := range o {
			v := v[k]
			args = append(args, v)
		}
		return args
	}(m)

	return fmt.Sprintf(`
        INSERT INTO "todo" (%s) 
        VALUES 
            (%s)
        RETURNING %s`, strings.Join(fs, ","), strings.Join(ph, ","), fields()), args
}

// UpdateQuery returns query statement and slice of arguments
func UpdateQuery(data *Todo) (string, []interface{}) {
	return fmt.Sprintf(`
        UPDATE "todo"
        SET 
            "title" = $1 , 
            "detail" = $2 , 
            "created_at" = $3 , 
            "updated_at" = $4 , 
            "deleted_at" = $5
        WHERE 
            "id" = $6
        RETURNING %s`, fields()),
		[]interface{}{data.Title, data.Detail, data.CreatedAt, data.UpdatedAt, data.DeletedAt, data.ID}
}

func deleteQuery(id int) (string, []interface{}) {
	now := time.Now().UTC()
	return `
        UPDATE "todo"
        SET  
            "deleted_at" = $1
        WHERE 
            "id" = $2`,
		[]interface{}{now, id}
}

func defaultFilter() string {
	return `"deleted_at" is NULL`
}

func scan(scanner rowScanner, data *Todo) error {
	var iID int
	var iTitle string
	var iDetail *string
	var iCreatedAt time.Time
	var iUpdatedAt time.Time
	var iDeletedAt **time.Time

	err := scanner.Scan(&iID, &iTitle, &iDetail, &iCreatedAt, &iUpdatedAt, &iDeletedAt)
	if err != nil {
		return err
	}

	data.ID = iID
	data.Title = iTitle
	if iDetail != nil {
		data.Detail = *iDetail
	}
	data.CreatedAt = iCreatedAt
	data.UpdatedAt = iUpdatedAt
	if iDeletedAt != nil {
		data.DeletedAt = *iDeletedAt
	}

	return nil
}

func scanRows(rows rowsResult) ([]*Todo, error) {
	collection := []*Todo{}
	for rows.Next() {
		data := &Todo{}
		err := scan(rows, data)
		if err != nil {
			return nil, err
		}
		collection = append(collection, data)
	}
	return collection, nil
}
