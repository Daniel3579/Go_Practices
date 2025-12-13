package integration

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"example.com/prc_integr/internal/db"
	"example.com/prc_integr/internal/httpapi"
	"example.com/prc_integr/internal/models"
	"example.com/prc_integr/internal/repo"
	"example.com/prc_integr/internal/service"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func withPostgresContainer(t *testing.T) (dsn string, cleanup func()) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Правильный API для testcontainers-go
	pgContainer, err := postgres.RunContainer(ctx,
		postgres.WithDatabase("notes_test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	require.NoError(t, err)

	host, err := pgContainer.Host(ctx)
	require.NoError(t, err)
	port, err := pgContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	dsn = fmt.Sprintf("postgres://test:test@%s:%s/notes_test?sslmode=disable",
		host, port.Port())

	cleanup = func() {
		_ = pgContainer.Terminate(ctx)
	}

	return dsn, cleanup
}

func newServerTC(t *testing.T, dsn string) *httptest.Server {
	t.Helper()

	dbx, err := sql.Open("postgres", dsn)
	require.NoError(t, err)

	err = dbx.Ping()
	require.NoError(t, err)

	db.MustApplyMigrations(dbx)
	err = db.TruncateAll(dbx)
	require.NoError(t, err)

	r := gin.Default()
	svc := service.Service{Notes: repo.NoteRepo{DB: dbx}}
	httpapi.Router{Svc: &svc}.Register(r)

	t.Cleanup(func() {
		_ = dbx.Close()
	})

	return httptest.NewServer(r)
}

func Test_CreateAndGet_TestContainers(t *testing.T) {
	dsn, cleanup := withPostgresContainer(t)
	defer cleanup()

	srv := newServerTC(t, dsn)
	defer srv.Close()

	// Create
	resp, err := http.Post(
		srv.URL+"/notes",
		"application/json",
		strings.NewReader(`{"title":"TC Test","content":"Testing with testcontainers"}`),
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var created models.Note
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	json.Unmarshal(body, &created)

	// Get
	resp2, err := http.Get(fmt.Sprintf("%s/notes/%d", srv.URL, created.ID))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp2.StatusCode)

	var retrieved models.Note
	body2, _ := io.ReadAll(resp2.Body)
	resp2.Body.Close()
	json.Unmarshal(body2, &retrieved)
	assert.Equal(t, created.ID, retrieved.ID)
}

func Test_CRUD_TestContainers(t *testing.T) {
	dsn, cleanup := withPostgresContainer(t)
	defer cleanup()

	srv := newServerTC(t, dsn)
	defer srv.Close()

	// CREATE
	resp, _ := http.Post(
		srv.URL+"/notes",
		"application/json",
		strings.NewReader(`{"title":"CRUD Test","content":"Full CRUD test"}`),
	)
	var note models.Note
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	json.Unmarshal(body, &note)
	assert.Greater(t, note.ID, int64(0))

	// READ
	resp2, _ := http.Get(fmt.Sprintf("%s/notes/%d", srv.URL, note.ID))
	assert.Equal(t, http.StatusOK, resp2.StatusCode)
	resp2.Body.Close()

	// UPDATE
	req, _ := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/notes/%d", srv.URL, note.ID),
		strings.NewReader(`{"title":"Updated CRUD","content":"Updated"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	resp3, _ := http.DefaultClient.Do(req)
	assert.Equal(t, http.StatusOK, resp3.StatusCode)
	resp3.Body.Close()

	// DELETE
	req2, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/notes/%d", srv.URL, note.ID), nil)
	resp4, _ := http.DefaultClient.Do(req2)
	assert.Equal(t, http.StatusOK, resp4.StatusCode)
	resp4.Body.Close()

	// Verify deleted
	resp5, _ := http.Get(fmt.Sprintf("%s/notes/%d", srv.URL, note.ID))
	assert.Equal(t, http.StatusNotFound, resp5.StatusCode)
	resp5.Body.Close()
}

func Test_Validation_TestContainers(t *testing.T) {
	dsn, cleanup := withPostgresContainer(t)
	defer cleanup()

	srv := newServerTC(t, dsn)
	defer srv.Close()

	// Empty title
	resp, _ := http.Post(
		srv.URL+"/notes",
		"application/json",
		strings.NewReader(`{"title":"","content":"No title"}`),
	)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	resp.Body.Close()

	// Invalid JSON
	resp2, _ := http.Post(
		srv.URL+"/notes",
		"application/json",
		strings.NewReader(`{invalid json}`),
	)
	assert.Equal(t, http.StatusBadRequest, resp2.StatusCode)
	resp2.Body.Close()
}
