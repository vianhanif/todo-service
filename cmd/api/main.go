package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/vianhanif/todo-service/internal/server/http"
)

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)

	s := http.NewServer()
	go s.Serve()
	log.Println("Service Started...")

	<-quit
	log.Println("Shutdown Server ...")
}
