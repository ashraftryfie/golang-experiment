package domain

import "errors"

var (
	ErrInvalidJobID       = errors.New("job id cannot be empty")
	ErrInvalidJobType     = errors.New("job type cannot be empty")
	ErrInvalidPayload     = errors.New("job payload cannot be empty")
	ErrJobNotFound        = errors.New("job not found")
	ErrJobAlreadyFinished = errors.New("job has already reached a terminal state")
	ErrMaxRetriesExceeded = errors.New("maximum retry attempts exceeded")
)
