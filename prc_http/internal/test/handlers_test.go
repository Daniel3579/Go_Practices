package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"example.com/prc_http/internal/api"
	"example.com/prc_http/internal/storage"
)

func setup() *api.Handlers {
	store := storage.NewMemoryStore()
	return api.NewHandlers(store)
}

func TestListTasks(t *testing.T) {
	h := setup()
	h.Store.Create("Task 1")
	h.Store.Create("Task 2")

	req, err := http.NewRequest("GET", "/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.ListTasks)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var tasks []*storage.Task
	if err := json.NewDecoder(rr.Body).Decode(&tasks); err != nil {
		t.Fatal(err)
	}

	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
}

func TestCreateTask(t *testing.T) {
	h := setup()
	task := map[string]string{"title": "New Task"}

	body, _ := json.Marshal(task)
	req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.CreateTask)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var createdTask storage.Task
	if err := json.NewDecoder(rr.Body).Decode(&createdTask); err != nil {
		t.Fatal(err)
	}

	if createdTask.Title != "New Task" {
		t.Errorf("expected task title to be 'New Task', got '%s'", createdTask.Title)
	}
}

func TestGetTask(t *testing.T) {
	h := setup()
	task := h.Store.Create("Get Task")

	req, err := http.NewRequest("GET", "/tasks/"+strconv.FormatInt(task.ID, 10), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetTask)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var retrievedTask storage.Task
	if err := json.NewDecoder(rr.Body).Decode(&retrievedTask); err != nil {
		t.Fatal(err)
	}

	if retrievedTask.ID != task.ID {
		t.Errorf("expected task ID to be %d, got %d", task.ID, retrievedTask.ID)
	}
}

func TestPatchTask(t *testing.T) {
	h := setup()
	task := h.Store.Create("Patch Task")

	req, err := http.NewRequest("PATCH", "/tasks/"+strconv.FormatInt(task.ID, 10), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.PatchTask)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var updatedTask storage.Task
	if err := json.NewDecoder(rr.Body).Decode(&updatedTask); err != nil {
		t.Fatal(err)
	}

	if !updatedTask.Done {
		t.Errorf("expected task to be marked as done")
	}
}

func TestDeleteTask(t *testing.T) {
	h := setup()
	task := h.Store.Create("Delete Task")

	req, err := http.NewRequest("DELETE", "/tasks/"+strconv.FormatInt(task.ID, 10), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.DeleteTask)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

	// Verify that the task is deleted
	_, err = h.Store.Get(task.ID)
	if err == nil {
		t.Errorf("expected task to be deleted, but it was found")
	}
}

func TestCreateTaskInvalidJSON(t *testing.T) {
	h := setup()
	invalidJSON := []byte(`{invalid json}`)

	req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(invalidJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.CreateTask)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestCreateTaskEmptyTitle(t *testing.T) {
	h := setup()
	task := map[string]string{"title": ""}

	body, _ := json.Marshal(task)
	req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.CreateTask)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestGetTaskNotFound(t *testing.T) {
	h := setup()

	req, err := http.NewRequest("GET", "/tasks/999", nil) // Non-existent ID
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetTask)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestPatchTaskNotFound(t *testing.T) {
	h := setup()

	req, err := http.NewRequest("PATCH", "/tasks/999", nil) // Non-existent ID
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.PatchTask)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestDeleteTaskNotFound(t *testing.T) {
	h := setup()

	req, err := http.NewRequest("DELETE", "/tasks/999", nil) // Non-existent ID
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.DeleteTask)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}
