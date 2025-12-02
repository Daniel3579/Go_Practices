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

func (h *Handler) ListNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(h.Repo.Notes())
}

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
