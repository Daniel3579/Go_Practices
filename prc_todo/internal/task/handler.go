package task

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	repo *Repo
}

func NewHandler(repo *Repo) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.list)          // GET /tasks
	r.Post("/", h.create)       // POST /tasks
	r.Get("/{id}", h.get)       // GET /tasks/{id}
	r.Put("/{id}", h.update)    // PUT /tasks/{id}
	r.Delete("/{id}", h.delete) // DELETE /tasks/{id}
	return r
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	tasks := h.repo.List()

	tasks = filter(tasks, r)
	tasks = paginate(tasks, r)

	if tasks == nil {
		httpError(w, http.StatusNoContent, "")
		return
	}

	writeJSON(w, http.StatusOK, tasks)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, bad := parseID(w, r)
	if bad {
		return
	}
	t, err := h.repo.Get(id)
	if err != nil {
		httpError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, t)
}

type createReq struct {
	Title string `json:"title"`
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req createReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Title == "" {
		httpError(w, http.StatusBadRequest, "invalid json: require non-empty title")
		return

	} else if len(req.Title) < 3 {
		httpError(w, http.StatusBadRequest, "the title length must be at least 3")
		return

	} else if len(req.Title) > 100 {
		httpError(w, http.StatusBadRequest, "the title length must be less than 100")
		return
	}

	t := h.repo.Create(req.Title)
	h.repo.SaveTasks("db.json")
	writeJSON(w, http.StatusCreated, t)
}

type updateReq struct {
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id, bad := parseID(w, r)
	if bad {
		return
	}
	var req updateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Title == "" {
		httpError(w, http.StatusBadRequest, "invalid json: require non-empty title")
		return
	}
	t, err := h.repo.Update(id, req.Title, req.Done)
	if err != nil {
		httpError(w, http.StatusNotFound, err.Error())
		return
	}
	h.repo.SaveTasks("db.json")
	writeJSON(w, http.StatusOK, t)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id, bad := parseID(w, r)
	if bad {
		return
	}
	if err := h.repo.Delete(id); err != nil {
		httpError(w, http.StatusNotFound, err.Error())
		return
	}
	h.repo.SaveTasks("db.json")
	w.WriteHeader(http.StatusNoContent)
}

// helpers

func parseID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	raw := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		httpError(w, http.StatusBadRequest, "invalid id")
		return 0, true
	}
	return id, false
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func filter(tasks []*Task, r *http.Request) []*Task {
	done := r.URL.Query().Get("done")
	if done != "" {
		filtered := tasks[:0]
		for _, t := range tasks {
			if boolToString(t.Done) == done {
				filtered = append(filtered, t)
			}
		}
		tasks = filtered
	}
	return tasks
}

func paginate(tasks []*Task, r *http.Request) []*Task {
	// Получаем параметры page и limit
	pageParam := r.URL.Query().Get("page")
	limitParam := r.URL.Query().Get("limit")

	// Преобразуем параметры в целые числа
	page, err := strconv.Atoi(pageParam)
	if err != nil || page <= 0 {
		page = 1 // Устанавливаем значение по умолчанию
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 5 // Устанавливаем значение по умолчанию
	}

	// Рассчитываем начальный индекс
	start := (page - 1) * limit
	if start >= len(tasks) {
		return nil
	}

	// Определяем конечный индекс
	end := start + limit
	if end > len(tasks) {
		end = len(tasks) // Ограничиваем, чтобы не выходить за пределы массива
	}

	// Получаем отфильтрованный список задач
	return tasks[start:end]
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func httpError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}
