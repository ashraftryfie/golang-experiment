package taskmanager

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Task struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`   // "pending" or "completed"
	Priority  int       `json:"priority"` // 1 to 5
	CreatedAt time.Time `json:"created_at"`
}

type CreateTaskPayload struct {
	Title    string `json:"title"`
	Priority int    `json:"priority"`
	Status   string `json:"status"`
}

type PaginationMeta struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

type PaginatedTasksResponse struct {
	Data       []Task         `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

type TaskStore struct {
	mu     sync.RWMutex
	tasks  map[string]Task
	order  []string
	seqGen int
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks: make(map[string]Task),
		order: make([]string, 0),
	}
}

func (s *TaskStore) Create(payload CreateTaskPayload) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.seqGen++
	id := fmt.Sprintf("task-%d", s.seqGen)

	status := payload.Status
	if status == "" {
		status = "pending"
	}

	task := Task{
		ID:        id,
		Title:     payload.Title,
		Status:    status,
		Priority:  payload.Priority,
		CreatedAt: time.Now().UTC(),
	}

	s.tasks[id] = task
	s.order = append(s.order, id)
	return task, nil
}

func (s *TaskStore) Get(id string) (Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tasks[id]
	return t, ok
}

func (s *TaskStore) Update(id string, payload CreateTaskPayload) (Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, ok := s.tasks[id]
	if !ok {
		return Task{}, false
	}

	existing.Title = payload.Title
	if payload.Status != "" {
		existing.Status = payload.Status
	}
	if payload.Priority > 0 {
		existing.Priority = payload.Priority
	}

	s.tasks[id] = existing
	return existing, true
}

func (s *TaskStore) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[id]; !ok {
		return false
	}
	delete(s.tasks, id)

	newOrder := make([]string, 0, len(s.order)-1)
	for _, k := range s.order {
		if k != id {
			newOrder = append(newOrder, k)
		}
	}
	s.order = newOrder
	return true
}

func (s *TaskStore) Query(statusFilter string, page, limit int) PaginatedTasksResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var filtered []Task
	for _, id := range s.order {
		t := s.tasks[id]
		if statusFilter == "" || t.Status == statusFilter {
			filtered = append(filtered, t)
		}
	}

	total := len(filtered)
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	start := (page - 1) * limit
	if start > total {
		start = total
	}
	end := start + limit
	if end > total {
		end = total
	}

	pageData := filtered[start:end]
	if pageData == nil {
		pageData = []Task{}
	}

	return PaginatedTasksResponse{
		Data: pageData,
		Pagination: PaginationMeta{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}
}

// NewTaskRouter wires endpoints for Task management.
func NewTaskRouter(store *TaskStore) http.Handler {
	mux := http.NewServeMux()

	// GET /tasks
	mux.HandleFunc("GET /tasks", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		status := q.Get("status")
		page, _ := strconv.Atoi(q.Get("page"))
		limit, _ := strconv.Atoi(q.Get("limit"))

		res := store.Query(status, page, limit)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(res)
	})

	// POST /tasks
	mux.HandleFunc("POST /tasks", func(w http.ResponseWriter, r *http.Request) {
		var payload CreateTaskPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
			return
		}

		if strings.TrimSpace(payload.Title) == "" {
			http.Error(w, `{"error":"title cannot be empty"}`, http.StatusUnprocessableEntity)
			return
		}
		if payload.Priority < 1 || payload.Priority > 5 {
			http.Error(w, `{"error":"priority must be between 1 and 5"}`, http.StatusUnprocessableEntity)
			return
		}

		task, _ := store.Create(payload)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", "/tasks/"+task.ID)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(task)
	})

	// GET /tasks/{id}
	mux.HandleFunc("GET /tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		task, ok := store.Get(id)
		if !ok {
			http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(task)
	})

	// PUT /tasks/{id}
	mux.HandleFunc("PUT /tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var payload CreateTaskPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
			return
		}

		task, ok := store.Update(id, payload)
		if !ok {
			http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(task)
	})

	// DELETE /tasks/{id}
	mux.HandleFunc("DELETE /tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if !store.Delete(id) {
			http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	return mux
}
