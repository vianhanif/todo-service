package app

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/vianhanif/todo-service/config"
	"github.com/vianhanif/todo-service/internal/alert"
	v1 "github.com/vianhanif/todo-service/internal/service/v1"
	"github.com/vianhanif/todo-service/internal/storage/todo"
)

// App .
type App struct {
	Config   *config.Config
	Alert    alert.Alert
	Storages *Storages
	Services *Services
	Clients  *Clients
}

// Storages .
type Storages struct {
	Todo todo.IStorage
}

// Services .
type Services struct {
	DB   *sql.DB
	Todo v1.Service
}

// Clients .
type Clients struct {
}

// NewApp .
func NewApp(
	config *config.Config,
	alertClient alert.Alert,
	storages *Storages,
	services *Services,
	clients *Clients,
) *App {

	return &App{
		Config:   config,
		Alert:    alertClient,
		Storages: storages,
		Services: services,
		Clients:  clients,
	}
}

type k string

const key = k("app")

// FromContext get app from context
func FromContext(ctx context.Context) *App {
	app, ok := ctx.Value(key).(*App)
	if !ok {
		return nil
	}
	return app
}

// ToContext put app to context
func ToContext(ctx context.Context, app *App) context.Context {
	return context.WithValue(ctx, key, app)
}

// InjectorMiddleware ...
func InjectorMiddleware(app *App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := ToContext(r.Context(), app)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
