package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

type TaskHandler struct {
	store TaskStore
	mux   *http.ServeMux
}

func NewTaskHandler(store TaskStore) *TaskHandler {
	h := &TaskHandler{
		store: store,
		mux:   http.NewServeMux(),
	}
	h.routes()
	return h
}

func (h *TaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *TaskHandler) routes() {
	h.mux.HandleFunc("POST /tasks", h.handleCreate)
	h.mux.HandleFunc("GET /tasks", h.handleList)
	h.mux.HandleFunc("GET /tasks/{id}", h.handleGet)
	h.mux.HandleFunc("PUT /tasks/{id}", h.handleUpdate)
	h.mux.HandleFunc("DELETE /tasks/{id}", h.handleDelete)
}

type createTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (h *TaskHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	task, err := h.store.Create(r.Context(), req.Title, req.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) handleList(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.store.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	task, err := h.store.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(task)
}

type updateTaskRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
}

func (h *TaskHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	task, err := h.store.Update(r.Context(), id, req.Title, req.Description, req.Status)
	if err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.store.Delete(r.Context(), id); err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	// Composition root for standalone execution
	store := NewInMemoryTaskStore()
	handler := NewTaskHandler(store)
	_ = http.ListenAndServe(":8080", handler)
}
