package main

import (
	"database/sql"
	"log"

	"example.com/prc_integr/internal/db"
	"example.com/prc_integr/internal/httpapi"
	"example.com/prc_integr/internal/repo"
	"example.com/prc_integr/internal/service"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	dsn := "postgres://test:test@localhost:54321/notes_test?sslmode=disable"

	dbx, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to connect to db:", err)
	}
	defer dbx.Close()

	db.MustApplyMigrations(dbx)

	r := gin.Default()
	svc := service.Service{Notes: repo.NoteRepo{DB: dbx}}
	httpapi.Router{Svc: &svc}.Register(r)

	log.Println("API listening on :8080")
	r.Run(":8080")
}
