package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"example.com/prc_http/internal/storage"
)

type Handlers struct {
	Store *storage.MemoryStore
}

func NewHandlers(store *storage.MemoryStore) *Handlers {
	return &Handlers{Store: store}
}

func parseIDFromURL(w http.ResponseWriter, r *http.Request) (int64, error) {
	// Разделение URL
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 2 {
		NotFound(w, "invalid path")
		return -1, errors.New("not found")
	}

	// Конвертация ID в int64
	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || id <= 0 {
		BadRequest(w, "invalid ID")
		return -1, errors.New("bad request")
	}

	return id, nil
}

// GET /tasks
func (h *Handlers) ListTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.Store.List()

	// Поддержка простых фильтров через query: ?q=text
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q != "" {
		filtered := tasks[:0]
		for _, t := range tasks {
			if strings.Contains(strings.ToLower(t.Title), strings.ToLower(q)) {
				filtered = append(filtered, t)
			}
		}
		tasks = filtered
	}

	JSON(w, http.StatusOK, tasks)
}

type createTaskRequest struct {
	Title string `json:"title"`
}

// POST /tasks
func (h *Handlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "" && !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		BadRequest(w, "Content-Type must be application/json")
		return
	}

	var req createTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "invalid json: "+err.Error())
		return
	}

	if req.Title == "" {
		BadRequest(w, "title is required")
		return

	} else if len(req.Title) > 140 {
		BadRequest(w, "invalid title length: must be between 1 and 140 characters")
		return
	}

	t := h.Store.Create(req.Title)
	JSON(w, http.StatusCreated, t)
}

// GET /tasks/{id}
func (h *Handlers) GetTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromURL(w, r)
	if err != nil {
		return
	}

	t, err := h.Store.Get(id)
	if err != nil {
		NotFound(w, "task not found")
		return
	}

	JSON(w, http.StatusOK, t)
}

// PATCH /tasks/{id}
func (h *Handlers) PatchTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromURL(w, r)
	if err != nil {
		return
	}

	t, err := h.Store.SetDone(id, true)
	if err != nil {
		NotFound(w, "task not found")
		return
	}

	JSON(w, http.StatusOK, t)
}

// DELETE /tasks/{id}
func (h *Handlers) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromURL(w, r)
	if err != nil {
		return
	}

	err = h.Store.Delete(id)
	if err != nil {
		NotFound(w, "task not found")
		return
	}

	JSON(w, http.StatusNoContent, nil)
}
