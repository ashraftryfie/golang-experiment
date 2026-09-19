package domain

import (
	"strings"
	"time"
)

type JobStatus string

const (
	StatusPending    JobStatus = "PENDING"
	StatusProcessing JobStatus = "PROCESSING"
	StatusCompleted  JobStatus = "COMPLETED"
	StatusFailed     JobStatus = "FAILED"
)

type Job struct {
	ID          string            `json:"id"`
	Type        string            `json:"type"`
	Payload     map[string]string `json:"payload"`
	Status      JobStatus         `json:"status"`
	Attempts    int               `json:"attempts"`
	MaxRetries  int               `json:"max_retries"`
	LastError   string            `json:"last_error,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	CompletedAt *time.Time        `json:"completed_at,omitempty"`
}

func NewJob(id, jobType string, payload map[string]string, maxRetries int) (*Job, error) {
	if strings.TrimSpace(id) == "" {
		return nil, ErrInvalidJobID
	}
	if strings.TrimSpace(jobType) == "" {
		return nil, ErrInvalidJobType
	}
	if len(payload) == 0 {
		return nil, ErrInvalidPayload
	}
	if maxRetries < 0 {
		maxRetries = 0
	}

	payloadCopy := make(map[string]string, len(payload))
	for k, v := range payload {
		payloadCopy[k] = v
	}

	now := time.Now().UTC()
	return &Job{
		ID:         id,
		Type:       jobType,
		Payload:    payloadCopy,
		Status:     StatusPending,
		Attempts:   0,
		MaxRetries: maxRetries,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

func (j *Job) MarkProcessing() error {
	if j.Status == StatusCompleted || j.Status == StatusFailed {
		return ErrJobAlreadyFinished
	}
	j.Status = StatusProcessing
	j.Attempts++
	j.UpdatedAt = time.Now().UTC()
	return nil
}

func (j *Job) MarkCompleted() error {
	if j.Status == StatusCompleted || j.Status == StatusFailed {
		return ErrJobAlreadyFinished
	}
	now := time.Now().UTC()
	j.Status = StatusCompleted
	j.UpdatedAt = now
	j.CompletedAt = &now
	j.LastError = ""
	return nil
}

func (j *Job) MarkFailed(errMsg string) {
	now := time.Now().UTC()
	j.UpdatedAt = now
	j.LastError = errMsg
	if j.Attempts > j.MaxRetries {
		j.Status = StatusFailed
		j.CompletedAt = &now
	} else {
		j.Status = StatusPending // Eligible for retry
	}
}

func (j *Job) IsTerminal() bool {
	return j.Status == StatusCompleted || j.Status == StatusFailed
}
