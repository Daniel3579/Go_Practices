package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"example.com/prc_mongo/internal/db"
	"example.com/prc_mongo/internal/notes"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI must be set")
	}

	dbName := os.Getenv("MONGO_DB")
	if dbName == "" {
		log.Fatal("MONGO_DB must be set")
	}

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		log.Fatal("HTTP_ADDR must be set")
	}

	deps, err := db.ConnectMongo(context.Background(), uri, dbName)
	if err != nil {
		log.Fatal("mongo connect:", err)
	}
	defer deps.Client.Disconnect(context.Background())

	repo, err := notes.NewRepo(deps.Database)
	if err != nil {
		log.Fatal("notes repo:", err)
	}
	h := notes.NewHandler(repo)

	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})
	r.Mount("/api/v1/notes", h.Routes())

	log.Println("listening on", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
