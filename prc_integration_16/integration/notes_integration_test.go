package integration

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"example.com/prc_integr/internal/db"
	"example.com/prc_integr/internal/httpapi"
	"example.com/prc_integr/internal/models"
	"example.com/prc_integr/internal/repo"
	"example.com/prc_integr/internal/service"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var dbDSN = flag.String("db-dsn", "", "Database connection string")

func newServer(t *testing.T, dsn string) *httptest.Server {
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

func Test_CreateAndGet_DockerCompose(t *testing.T) {
	dsn := *dbDSN
	if dsn == "" {
		dsn = os.Getenv("DB_DSN")
	}
	if dsn == "" {
		t.Skip("DB_DSN not set. Run with -db-dsn flag or set DB_DSN env var")
	}

	srv := newServer(t, dsn)
	defer srv.Close()

	// Create note
	resp, err := http.Post(
		srv.URL+"/notes",
		"application/json",
		strings.NewReader(`{"title":"Integration Test","content":"Testing with docker-compose"}`),
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	var created models.Note
	err = json.Unmarshal(body, &created)
	require.NoError(t, err)
	assert.Equal(t, "Integration Test", created.Title)
	assert.Greater(t, created.ID, int64(0))

	// Get note
	resp2, err := http.Get(fmt.Sprintf("%s/notes/%d", srv.URL, created.ID))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp2.StatusCode)

	body2, _ := io.ReadAll(resp2.Body)
	resp2.Body.Close()

	var retrieved models.Note
	err = json.Unmarshal(body2, &retrieved)
	require.NoError(t, err)
	assert.Equal(t, created.ID, retrieved.ID)
	assert.Equal(t, "Integration Test", retrieved.Title)
}

func Test_Update_DockerCompose(t *testing.T) {
	dsn := *dbDSN
	if dsn == "" {
		dsn = os.Getenv("DB_DSN")
	}
	if dsn == "" {
		t.Skip("DB_DSN not set. Run with -db-dsn flag or set DB_DSN env var")
	}

	srv := newServer(t, dsn)
	defer srv.Close()

	// Create
	resp, _ := http.Post(
		srv.URL+"/notes",
		"application/json",
		strings.NewReader(`{"title":"Original","content":"Original content"}`),
	)
	var created models.Note
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	json.Unmarshal(body, &created)

	// Update
	req, _ := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/notes/%d", srv.URL, created.ID),
		strings.NewReader(`{"title":"Updated","content":"Updated content"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	resp2, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp2.StatusCode)
	resp2.Body.Close()

	// Verify
	resp3, _ := http.Get(fmt.Sprintf("%s/notes/%d", srv.URL, created.ID))
	var updated models.Note
	body3, _ := io.ReadAll(resp3.Body)
	resp3.Body.Close()
	json.Unmarshal(body3, &updated)
	assert.Equal(t, "Updated", updated.Title)
}

func Test_Delete_DockerCompose(t *testing.T) {
	dsn := *dbDSN
	if dsn == "" {
		dsn = os.Getenv("DB_DSN")
	}
	if dsn == "" {
		t.Skip("DB_DSN not set. Run with -db-dsn flag or set DB_DSN env var")
	}

	srv := newServer(t, dsn)
	defer srv.Close()

	// Create
	resp, _ := http.Post(
		srv.URL+"/notes",
		"application/json",
		strings.NewReader(`{"title":"To Delete","content":"This will be deleted"}`),
	)
	var created models.Note
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	json.Unmarshal(body, &created)

	// Delete
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/notes/%d", srv.URL, created.ID), nil)
	resp2, _ := http.DefaultClient.Do(req)
	require.Equal(t, http.StatusOK, resp2.StatusCode)
	resp2.Body.Close()

	// Verify it's gone
	resp3, _ := http.Get(fmt.Sprintf("%s/notes/%d", srv.URL, created.ID))
	require.Equal(t, http.StatusNotFound, resp3.StatusCode)
	resp3.Body.Close()
}

func Test_List_DockerCompose(t *testing.T) {
	dsn := *dbDSN
	if dsn == "" {
		dsn = os.Getenv("DB_DSN")
	}
	if dsn == "" {
		t.Skip("DB_DSN not set. Run with -db-dsn flag or set DB_DSN env var")
	}

	srv := newServer(t, dsn)
	defer srv.Close()

	// Create multiple notes
	for i := 1; i <= 3; i++ {
		http.Post(
			srv.URL+"/notes",
			"application/json",
			strings.NewReader(fmt.Sprintf(`{"title":"Note %d","content":"Content %d"}`, i, i)),
		)
	}

	// List
	resp, _ := http.Get(srv.URL + "/notes?limit=10&offset=0")
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string][]models.Note
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	json.Unmarshal(body, &result)
	assert.Len(t, result["data"], 3)
}

func Test_NotFound_DockerCompose(t *testing.T) {
	dsn := *dbDSN
	if dsn == "" {
		dsn = os.Getenv("DB_DSN")
	}
	if dsn == "" {
		t.Skip("DB_DSN not set. Run with -db-dsn flag or set DB_DSN env var")
	}

	srv := newServer(t, dsn)
	defer srv.Close()

	resp, _ := http.Get(srv.URL + "/notes/99999")
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	resp.Body.Close()
}
