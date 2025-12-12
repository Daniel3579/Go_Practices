package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/prc_optim/internal/db"
	apihttp "example.com/prc_optim/internal/http"
	"example.com/prc_optim/internal/http/handlers"
	"example.com/prc_optim/internal/repo"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv"
)

func main() {

	var repository repo.NoteRepository

	databaseURL := os.Getenv("POSTGRES_URL_GO")
	if databaseURL == "" {
		databaseURL = "postgres://postgres@localhost:5433/prc_db?sslmode=disable"
	}

	sqlDB, err := db.NewDB(databaseURL, &db.Config{
		MaxOpenConns:    20,
		MaxIdleConns:    10,
		ConnMaxLifetime: 30 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
	})
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	repository = repo.NewNoteRepoPostgres(sqlDB)
	log.Println("✅ Using PostgreSQL")

	defer sqlDB.Close()

	router := gin.Default()
	h := &handlers.Handler{Repo: repository}

	// Вызов SetupRoutes из правильного пакета (алиас apihttp)
	apihttp.SetupRoutes(router, h)

	// Graceful shutdown
	server := &http.Server{Addr: ":8080", Handler: router}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Server starting on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-sigChan
	log.Println("Shutdown...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}
