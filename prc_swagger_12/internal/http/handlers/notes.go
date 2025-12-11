package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"example.com/prc_notes_api/internal/core"
	"example.com/prc_notes_api/internal/repo"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Repo *repo.NoteRepoMem
}

// CreateNote создает новую заметку
// @Summary      Создать заметку
// @Description  Создает новую заметку с указанным заголовком и содержимым
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        input  body      core.ReqNote  true  "Данные новой заметки"
// @Success      201    {object}  core.Note
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /notes [post]
func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var n core.ReqNote
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	create, err := h.Repo.Create(n)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(create)
}

// GetNote получает заметку по ID
// @Summary      Получить заметку
// @Description  Возвращает заметку по указанному ID
// @Tags         notes
// @Produce      json
// @Param        id   path      int  true  "ID заметки"
// @Success      200  {object}  core.Note
// @Failure      404  {object}  map[string]string
// @Router       /notes/{id} [get]
func (h *Handler) GetNote(w http.ResponseWriter, r *http.Request) {
	id, bad := parseID(w, r)
	if bad {
		return
	}
	note, err := h.Repo.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(note)
}

// ListNotes возвращает список всех заметок
// @Summary      Список заметок
// @Description  Возвращает список всех заметок
// @Tags         notes
// @Produce      json
// @Success      200    {array}  core.Note
// @Failure      500    {object}  map[string]string
// @Router       /notes [get]
func (h *Handler) ListNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(h.Repo.Notes())
}

// Delete удаляет заметку по ID
// @Summary      Удалить заметку
// @Description  Удаляет заметку по указанному ID
// @Tags         notes
// @Param        id  path  int  true  "ID заметки"
// @Success      204  "No Content"
// @Failure      404  {object}  map[string]string
// @Router       /notes/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, bad := parseID(w, r)
	if bad {
		return
	}
	err := h.Repo.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// Update обновляет заметку
// @Summary      Обновить заметку
// @Description  Обновляет заметку по указанному ID
// @Tags         notes
// @Accept       json
// @Param        id     path   int        true  "ID заметки"
// @Param        input  body   core.ReqNote true  "Поля для обновления"
// @Success      204    "No Content"
// @Failure      400    {object}  map[string]string
// @Failure      404    {object}  map[string]string
// @Router       /notes/{id} [patch]
// @Router       /notes/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, bad := parseID(w, r)
	if bad {
		return
	}
	var n core.ReqNote
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
	}

	h.Repo.Update(id, &n)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func parseID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	raw := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(raw, 10, 64)
	bad := false
	if err != nil || id <= 0 {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(422)
		id = 0
		bad = true
	}
	return id, bad
}
