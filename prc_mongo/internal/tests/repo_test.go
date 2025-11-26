package notes

import (
	"context"
	"log"
	"os"
	"testing"

	"example.com/prc_mongo/internal/db"
	"example.com/prc_mongo/internal/notes"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func TestCreateAndGet(t *testing.T) {
	ctx := context.Background()
	deps, err := db.ConnectMongo(ctx, os.Getenv("MONGO_URI"), "pz8_test")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := deps.Client.Database("pz8_test").Drop(ctx); err != nil {
			t.Error(err)
		}
		deps.Client.Disconnect(ctx)
	})

	r, err := notes.NewRepo(deps.Database)
	if err != nil {
		t.Fatal(err)
	}

	created, err := r.Create(ctx, "T1", "C1", nil)
	if err != nil {
		t.Fatal(err)
	}

	got, err := r.ByID(ctx, created.ID.Hex())
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != "T1" {
		t.Fatalf("want T1 got %s", got.Title)
	}
}
