package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	ErrNoteNotFound    = errors.New("note not found")
	ErrEmptyTitle      = errors.New("title cannot be empty")
	ErrInvalidNoteID   = errors.New("invalid note id")
)

type Note struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
}

type NoteStore struct {
	filePath string
	mu       sync.RWMutex
	notes    []Note
	nextID   int
}

func NewNoteStore(filePath string) (*NoteStore, error) {
	store := &NoteStore{
		filePath: filePath,
		notes:    make([]Note, 0),
		nextID:   1,
	}

	if err := store.load(); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *NoteStore) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("reading notes file: %w", err)
	}

	if len(data) == 0 {
		return nil
	}

	if err := json.Unmarshal(data, &s.notes); err != nil {
		return fmt.Errorf("parsing notes JSON: %w", err)
	}

	maxID := 0
	for _, n := range s.notes {
		if n.ID > maxID {
			maxID = n.ID
		}
	}
	s.nextID = maxID + 1
	return nil
}

func (s *NoteStore) save() error {
	data, err := json.MarshalIndent(s.notes, "", "  ")
	if err != nil {
		return fmt.Errorf("serializing notes: %w", err)
	}
	return os.WriteFile(s.filePath, data, 0644)
}

func (s *NoteStore) Add(title, content string, tags []string) (*Note, error) {
	if strings.TrimSpace(title) == "" {
		return nil, ErrEmptyTitle
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	note := Note{
		ID:        s.nextID,
		Title:     title,
		Content:   content,
		Tags:      tags,
		CreatedAt: time.Now().UTC(),
	}
	s.nextID++
	s.notes = append(s.notes, note)

	if err := s.save(); err != nil {
		return nil, err
	}
	return &note, nil
}

func (s *NoteStore) List() []Note {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]Note, len(s.notes))
	copy(result, s.notes)
	return result
}

func (s *NoteStore) Search(query string) []Note {
	s.mu.RLock()
	defer s.mu.RUnlock()

	q := strings.ToLower(query)
	var matches []Note
	for _, n := range s.notes {
		if strings.Contains(strings.ToLower(n.Title), q) ||
			strings.Contains(strings.ToLower(n.Content), q) {
			matches = append(matches, n)
			continue
		}
		for _, tag := range n.Tags {
			if strings.Contains(strings.ToLower(tag), q) {
				matches = append(matches, n)
				break
			}
		}
	}
	return matches
}

func (s *NoteStore) Delete(id int) error {
	if id <= 0 {
		return ErrInvalidNoteID
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	idx := -1
	for i, n := range s.notes {
		if n.ID == id {
			idx = i
			break
		}
	}
	if idx == -1 {
		return ErrNoteNotFound
	}

	s.notes = append(s.notes[:idx], s.notes[idx+1:]...)
	return s.save()
}
