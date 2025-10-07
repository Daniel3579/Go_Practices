package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"example.com/prc_todo/internal/task"
	myMW "example.com/prc_todo/pkg/middleware"
)

func main() {
	repo := task.NewRepo()
	h := task.NewHandler(repo)

	// Загружаем задачи при старте
	if err := repo.LoadTasks("db.json"); err != nil && !os.IsNotExist(err) {
		log.Fatalf("Error loading tasks: %v", err)
	}

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.Recoverer)
	r.Use(myMW.Logger)
	r.Use(myMW.SimpleCORS)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Route("/api/v1", func(api chi.Router) {
		api.Mount("/tasks", h.Routes())
	})

	addr := ":8080"
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
