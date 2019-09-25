// template version: 1.0.9
package todo_test

import (
	"context"
	s "database/sql"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/vianhanif/go-pkg/generator"
	queryable "github.com/vianhanif/todo-service/internal/storage"
	"github.com/vianhanif/todo-service/internal/storage/todo"

	_ "github.com/lib/pq"
)

func TestCreate(t *testing.T) {

	data, err := fakeTodo()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating fake data, data : %+v", err, data)
	}
	db, err := getDB()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbCtx := queryable.NewContext(context.TODO(), db)

	storage := todo.NewStorage(db)
	err = storage.Create(dbCtx, data)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when create data", err)
	}
}
func TestCreateUsingTrx(t *testing.T) {

	data, err := fakeTodo()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating fake data", err)
	}
	db, err := getDB()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when beginning transaction", err)
	}
	defer db.Close()
	txCtx := queryable.NewContext(context.TODO(), tx)

	storage := todo.NewStorage(tx)
	err = storage.Create(txCtx, data)
	if err != nil {
		tx.Rollback()
		t.Fatalf("an error '%s' was not expected when create data", err)
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		t.Fatalf("an error '%s' was not expected when create data", err)
	}
	if err != nil {
		tx.Rollback()
		t.Fatalf("an error '%s' was not expected when committing transaction", err)
	}
}
func TestUpdate(t *testing.T) {
	target, err := fakeCreate()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating fake data", err)
	}
	updateSource, err := fakeTodo()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating fake data", err)
	}
	db, err := getDB()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	dbCtx := queryable.NewContext(context.TODO(), db)
	target.CreatedAt = updateSource.CreatedAt
	target.DeletedAt = updateSource.DeletedAt
	target.Detail = updateSource.Detail
	target.Title = updateSource.Title
	target.UpdatedAt = updateSource.UpdatedAt

	storage := todo.NewStorage(db)
	err = storage.Update(dbCtx, target)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when updating data", err)
	}
}
func TestUpdateUsingTrx(t *testing.T) {

	target, err := fakeCreate()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating fake data", err)
	}
	updateSource, err := fakeTodo()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating fake data", err)
	}
	db, err := getDB()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when beginning transaction", err)
	}
	defer db.Close()
	txCtx := queryable.NewContext(context.TODO(), tx)
	target.CreatedAt = updateSource.CreatedAt
	target.DeletedAt = updateSource.DeletedAt
	target.Detail = updateSource.Detail
	target.Title = updateSource.Title
	target.UpdatedAt = updateSource.UpdatedAt

	storage := todo.NewStorage(tx)
	err = storage.Update(txCtx, target)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when updating data", err)
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		t.Fatalf("an error '%s' was not expected when committing transaction", err)
	}
}
func TestSaveCreate(t *testing.T) {
	data, err := fakeTodo()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating fake data", err)
	}
	db, err := getDB()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	dbCtx := queryable.NewContext(context.TODO(), db)

	storage := todo.NewStorage(db)
	err = storage.Save(dbCtx, data)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when create data", err)
	}
}
func TestSaveCreateUsingTrx(t *testing.T) {

	data, err := fakeTodo()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating fake data", err)
	}
	db, err := getDB()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when beginning transaction", err)
	}
	defer db.Close()
	txCtx := queryable.NewContext(context.TODO(), tx)

	storage := todo.NewStorage(tx)
	err = storage.Save(txCtx, data)
	if err != nil {
		tx.Rollback()
		t.Fatalf("an error '%s' was not expected when save data", err)
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		t.Fatalf("an error '%s' was not expected when saving data", err)
	}
	if err != nil {
		tx.Rollback()
		t.Fatalf("an error '%s' was not expected when committing transaction", err)
	}
}
func TestSaveUpdate(t *testing.T) {
	target, err := fakeCreate()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating fake data", err)
	}
	updateSource, err := fakeTodo()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating fake data", err)
	}
	db, err := getDB()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	dbCtx := queryable.NewContext(context.TODO(), db)
	target.CreatedAt = updateSource.CreatedAt
	target.DeletedAt = updateSource.DeletedAt
	target.Detail = updateSource.Detail
	target.Title = updateSource.Title
	target.UpdatedAt = updateSource.UpdatedAt

	storage := todo.NewStorage(db)
	err = storage.Save(dbCtx, target)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when saving data", err)
	}
}
func TestSaveUpdateUsingTrx(t *testing.T) {
	target, err := fakeCreate()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating fake data", err)
	}
	updateSource, err := fakeTodo()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating fake data", err)
	}
	db, err := getDB()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when beginning transaction", err)
	}
	defer db.Close()
	txCtx := queryable.NewContext(context.TODO(), tx)
	target.CreatedAt = updateSource.CreatedAt
	target.DeletedAt = updateSource.DeletedAt
	target.Detail = updateSource.Detail
	target.Title = updateSource.Title
	target.UpdatedAt = updateSource.UpdatedAt

	storage := todo.NewStorage(tx)
	err = storage.Save(txCtx, target)
	if err != nil {
		tx.Rollback()
		t.Fatalf("an error '%s' was not expected when saving data", err)
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		t.Fatalf("an error '%s' was not expected when committing transaction", err)
	}
}
func TestDelete(t *testing.T) {
	target, err := fakeCreate()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating fake data", err)
	}
	db, err := getDB()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	dbCtx := queryable.NewContext(context.TODO(), db)

	storage := todo.NewStorage(db)
	err = storage.Delete(dbCtx,
		target.ID,
	)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when deleting data", err)
	}

	deletedData, err := storage.FindByKeys(dbCtx,
		target.ID,
	)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when finding deleted data", err)
	}
	if deletedData != nil {
		t.Fatalf("deleted data should not be returned when calling FindByKeys")
	}

	deletedData, err = storage.FindByKeysNoFilter(dbCtx,
		target.ID,
	)
	if deletedData == nil {
		t.Fatalf("deleted data should be returned when calling FindByKeysNoFilter")
	}
	if deletedData.DeletedAt == nil {
		t.Fatalf("deleted data should have valid 'DeletedAt' value")
	}
}

func TestDeleteUsingTrx(t *testing.T) {
	target, err := fakeCreate()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating fake data", err)
	}
	db, err := getDB()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when beginning transaction", err)
	}
	defer db.Close()
	txCtx := queryable.NewContext(context.TODO(), tx)
	dbCtx := queryable.NewContext(context.TODO(), db)

	storage := todo.NewStorage(tx)
	err = storage.Delete(txCtx,
		target.ID,
	)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when deleting data", err)
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		t.Fatalf("an error '%s' was not expected when committing transaction", err)
	}
	storage = todo.NewStorage(db)
	deletedData, err := storage.FindByKeys(dbCtx,
		target.ID,
	)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when finding deleted data", err)
	}
	if deletedData != nil {
		t.Fatalf("deleted data should not be returned when calling FindByKeys")
	}

	deletedData, err = storage.FindByKeysNoFilter(dbCtx,
		target.ID,
	)
	if deletedData == nil {
		t.Fatalf("deleted data should be returned when calling FindByKeysNoFilter")
	}
	if deletedData.DeletedAt == nil {
		t.Fatalf("deleted data should have valid 'DeletedAt' value")
	}
}

func TestSingle(t *testing.T) {

	target, err := fakeCreate()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating fake data", err)
	}
	db, err := getDB()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	dbCtx := queryable.NewContext(context.TODO(), db)

	storage := todo.NewStorage(db)
	data, err := storage.Single(dbCtx, `"id" = $1`,
		target.ID,
	)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when finding with single method", err)
	}
	if data == nil {
		t.Fatalf("undeleted data should be returned when calling Single")
	}
}
func TestFirst(t *testing.T) {
	target, err := fakeCreate()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating fake data", err)
	}
	db, err := getDB()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	dbCtx := queryable.NewContext(context.TODO(), db)

	storage := todo.NewStorage(db)
	data, err := storage.First(dbCtx, `"id" = $1`,
		target.ID,
	)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when finding with first method", err)
	}
	if data == nil {
		t.Fatalf("undeleted data should be returned when calling First")
	}
}
func TestFirstOrder(t *testing.T) {

	target, err := fakeCreate()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating fake data", err)
	}
	db, err := getDB()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	dbCtx := queryable.NewContext(context.TODO(), db)

	storage := todo.NewStorage(db)
	data, err := storage.FirstOrder(dbCtx, `"id" = $1`, `"id" asc`,
		target.ID,
	)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when finding with first method", err)
	}
	if data == nil {
		t.Fatalf("undeleted data should be returned when calling First")
	}
}
func TestFindAll(t *testing.T) {
	_, err := fakeCreate()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating fake data", err)
	}
	db, err := getDB()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbCtx := queryable.NewContext(context.TODO(), db)

	storage := todo.NewStorage(db)
	result, err := storage.FindAll(dbCtx, 1, 10, "created_at")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when querying data", err)
	}
	if len(result) < 1 {
		t.Fatalf("query result expected has length greater than 0.")
	}
	for _, row := range result {
		if row.DeletedAt != nil {
			t.Fatalf("logically deleted data is not expected to be included in query result")
		}
	}

}
func TestFindByKeys(t *testing.T) {

	target, err := fakeCreate()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating fake data", err)
	}
	db, err := getDB()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbCtx := queryable.NewContext(context.TODO(), db)

	storage := todo.NewStorage(db)
	data, err := storage.FindByKeys(dbCtx,
		target.ID,
	)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when finding data by keys", err)
	}
	if data == nil {
		t.Fatalf("undeleted data should be returned when calling FindByKeys")
	}
}
func TestWhere(t *testing.T) {

	target, err := fakeCreate()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating fake data", err)
	}
	db, err := getDB()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbCtx := queryable.NewContext(context.TODO(), db)

	storage := todo.NewStorage(db)
	result, err := storage.Where(dbCtx, "\"id\" = $1", target.ID)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when querying data", err)
	}
	if len(result) < 1 {
		t.Fatalf("query result expected has length greater than 0.")
	}
}
func TestWhereOrder(t *testing.T) {
	target, err := fakeCreate()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating fake data", err)
	}
	db, err := getDB()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbCtx := queryable.NewContext(context.TODO(), db)

	storage := todo.NewStorage(db)
	result, err := storage.WhereOrder(dbCtx, "\"id\" = $1", "\"id\" asc", target.ID)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when querying data", err)
	}
	if len(result) < 1 {
		t.Fatalf("query result expected has length greater than 0.")
	}
}
func TestWhereWithPaging(t *testing.T) {
	target, err := fakeCreate()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when generating fake data", err)
	}
	db, err := getDB()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbCtx := queryable.NewContext(context.TODO(), db)

	storage := todo.NewStorage(db)
	result, err := storage.WhereWithPaging(dbCtx, 1, 1, "\"id\" = $1", "\"id\" asc", target.ID)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when querying data", err)
	}
	if len(result) < 1 {
		t.Fatalf("query result expected has length greater than 0.")
	}
}
func TestWhereNoFilter(t *testing.T) {

	target, err := fakeCreate()
	if err != nil {
		t.Fatalf("an error '%s' was not expecting when generating fake data", err)
	}
	db, err := getDB()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbCtx := queryable.NewContext(context.TODO(), db)

	storage := todo.NewStorage(db)
	result, err := storage.WhereNoFilter(dbCtx, "\"id\" = $1", target.ID)
	if err != nil {
		t.Fatalf("an error '%s' was not expecting when querying data", err)
	}
	if len(result) < 1 {
		t.Fatalf("query result expecting has length greater than 0.")
	}
}
func fakeTodo() (*todo.Todo, error) {
	fake := &todo.Todo{}

	randID, _ := generator.RandomNumericString(8)
	i64randID, _ := strconv.ParseInt(randID, 10, 64)
	irandID := int(i64randID)
	fake.ID = irandID

	randTitle, _ := generator.RandomStringSet(255, "abcdefghijklmnopqrstuvwxyz")
	fake.Title = randTitle

	randDetail, _ := generator.RandomStringSet(32, "abcdefghijklmnopqrstuvwxyz")
	fake.Detail = randDetail

	nowCreatedAt := time.Now()
	fake.CreatedAt = nowCreatedAt

	nowUpdatedAt := time.Now()
	fake.UpdatedAt = nowUpdatedAt

	fake.DeletedAt = nil
	return fake, nil
}

// fakeCreate , create fake data used to simplify data creation on test functions.
func fakeCreate() (*todo.Todo, error) {

	data, err := fakeTodo()
	if err != nil {
		return nil, err
	}

	db, err := getDB()
	if err != nil {
		return nil, err
	}
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	txCtx := queryable.NewContext(context.TODO(), tx)

	storage := todo.NewStorage(tx)
	err = storage.Create(txCtx, data)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return data, nil
}

func getDB() (*s.DB, error) {
	return s.Open("postgres", os.Getenv("TODO_DB"))
}
