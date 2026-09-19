package main

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"
)

var (
	ErrTaskNotFound = errors.New("task not found")
	ErrEmptyTitle   = errors.New("task title cannot be empty")
)

type TaskStatus string

const (
	StatusTodo       TaskStatus = "TODO"
	StatusInProgress TaskStatus = "IN_PROGRESS"
	StatusDone       TaskStatus = "DONE"
)

type Task struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type TaskStore interface {
	Create(ctx context.Context, title, desc string) (*Task, error)
	Get(ctx context.Context, id string) (*Task, error)
	List(ctx context.Context) ([]*Task, error)
	Update(ctx context.Context, id string, title, desc string, status TaskStatus) (*Task, error)
	Delete(ctx context.Context, id string) error
}

type InMemoryTaskStore struct {
	mu     sync.RWMutex
	tasks  map[string]*Task
	nextID int
}

func NewInMemoryTaskStore() *InMemoryTaskStore {
	return &InMemoryTaskStore{
		tasks:  make(map[string]*Task),
		nextID: 1,
	}
}

func (s *InMemoryTaskStore) Create(_ context.Context, title, desc string) (*Task, error) {
	if strings.TrimSpace(title) == "" {
		return nil, ErrEmptyTitle
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	task := &Task{
		ID:          strings.ToLower(strings.ReplaceAll(title, " ", "-")),
		Title:       title,
		Description: desc,
		Status:      StatusTodo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	s.tasks[task.ID] = task
	return cloneTask(task), nil
}

func (s *InMemoryTaskStore) Get(_ context.Context, id string) (*Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[id]
	if !exists {
		return nil, ErrTaskNotFound
	}
	return cloneTask(task), nil
}

func (s *InMemoryTaskStore) List(_ context.Context) ([]*Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*Task
	for _, t := range s.tasks {
		result = append(result, cloneTask(t))
	}
	return result, nil
}

func (s *InMemoryTaskStore) Update(_ context.Context, id string, title, desc string, status TaskStatus) (*Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[id]
	if !exists {
		return nil, ErrTaskNotFound
	}

	if title != "" {
		task.Title = title
	}
	if desc != "" {
		task.Description = desc
	}
	if status != "" {
		task.Status = status
	}
	task.UpdatedAt = time.Now().UTC()

	return cloneTask(task), nil
}

func (s *InMemoryTaskStore) Delete(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[id]; !exists {
		return ErrTaskNotFound
	}
	delete(s.tasks, id)
	return nil
}

func cloneTask(t *Task) *Task {
	c := *t
	return &c
}
