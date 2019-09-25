package test

import (
	"context"
	"database/sql"
	"sync"

	"github.com/icrowley/fake"

	// init
	_ "github.com/lib/pq"

	"github.com/vianhanif/todo-service/config"
	"github.com/vianhanif/todo-service/internal/alert"
	"github.com/vianhanif/todo-service/internal/app"
	v1 "github.com/vianhanif/todo-service/internal/service/v1"
	"github.com/vianhanif/todo-service/internal/storage"
	"github.com/vianhanif/todo-service/internal/storage/todo"
)

var (
	oApp     sync.Once
	storages *app.Storages
	services *app.Services
	clients  *app.Clients
	a        *app.App
)

// GetApp .
func GetApp(ctx context.Context) *app.App {
	oApp.Do(func() {
		if a == nil {
			cfg, err := config.GetConfiguration()
			if err != nil {
				panic(err)
			}
			db := getDB(cfg)
			storages := getStorages(db)
			services := getServices(db, storages)
			clients := getClients(cfg)
			a = &app.App{
				Config:   cfg,
				Alert:    alert.NewMockAlert(),
				Storages: storages,
				Services: services,
				Clients:  clients,
			}
		}
	})
	return a
}

func getStorages(db storage.Queryable) *app.Storages {
	return &app.Storages{
		Todo: todo.NewStorage(db),
	}
}

func getServices(db *sql.DB, stg *app.Storages) *app.Services {
	return &app.Services{
		DB:   db,
		Todo: v1.NewService(db, stg.Todo),
	}
}

func getClients(cfg *config.Config) *app.Clients {
	return &app.Clients{}
}

func getDB(cfg *config.Config) *sql.DB {
	db, err := sql.Open("postgres", cfg.DB)
	if err != nil {
		panic(err)
	}
	return db
}

// FakeTodo .
func FakeTodo() *todo.Todo {
	return &todo.Todo{
		Title:  fake.JobTitle(),
		Detail: fake.Sentence(),
	}
}
