package http

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/vianhanif/todo-service/internal/alert"
	v1 "github.com/vianhanif/todo-service/internal/service/v1"
	"github.com/vianhanif/todo-service/internal/storage/todo"

	// init
	_ "github.com/lib/pq"

	"github.com/vianhanif/todo-service/config"
	"github.com/vianhanif/todo-service/internal/app"
)

// Server represents the http server that handles the requests
type Server struct {
	app *app.App
}

// Serve serves http requests
func (hs *Server) Serve() {
	r := hs.compileRouter(hs.app.Config)

	port := hs.app.Config.HTTPPort
	log.Printf("About to listen on %s. Go to http://127.0.0.1:%s", port, port)
	srv := http.Server{Addr: ":" + port, Handler: r}

	// start server and gracefull exit when signal triggered
	sSignalChan := make(chan os.Signal, 1)
	signal.Notify(sSignalChan, os.Interrupt)

	sErrChan := make(chan error)
	go func() {
		sErrChan <- srv.ListenAndServe()
	}()

	select {
	case err := <-sSignalChan:
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
		if err != nil {
			log.Println("Error from sSignalChan")
			log.Println(err)
		}
	case err := <-sErrChan:
		log.Println("Error from sErrChan")
		log.Println(err)
	}

	log.Println("Signal off called from Channel, done exit")
}

// NewServer creates a new http server
func NewServer() *Server {
	config, err := config.GetConfiguration()
	if err != nil {
		panic(err)
	}
	db, err := sql.Open("postgres", config.DB)
	if err != nil {
		panic(err)
	}

	todoStorage := todo.NewStorage(db)
	todoService := v1.NewService(db, todoStorage)

	alertClient := alert.NewSlackAlert(alert.SlackAlertConfig{
		Token:   config.SlackToken,
		Channel: config.SlackAlertChannel,
	})

	storages := &app.Storages{
		Todo: todoStorage,
	}
	services := &app.Services{
		DB:   db,
		Todo: todoService,
	}

	a := app.NewApp(config, alertClient, storages, services, nil)

	return &Server{
		app: a,
	}
}
