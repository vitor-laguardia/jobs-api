package job

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type JobStatus string

const (
	JobStatusPending JobStatus = "pending"
	JobStatusRunning JobStatus = "running"
	JobStatusDone    JobStatus = "done"
	JobStatusFailed  JobStatus = "failed"
)

type JobPriority int

const (
	JobPriorityLow    JobPriority = 1
	JobPriorityMedium JobPriority = 2
	JobPriorityHigh   JobPriority = 3
)

type Job struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Status      JobStatus   `json:"status"`
	Priority    JobPriority `json:"priority"`
	UserID      string      `json:"user_id"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

const (
	MsgInvalidStatusTransition = "invalid job status transition"
	MsgJobNotFound             = "job not found"
	MsgUserNotFound            = "user not found"
	MsgPriorityOutOfRange      = "priority must be in range [1-3]"
	MsgInvalidStatus           = "invalid status value"
	MsgJobAlreadyExists        = "job already exists"
)

var (
	ErrJobNotFound         = errors.New(MsgJobNotFound)
	ErrJobStatusTransition = errors.New(MsgInvalidStatusTransition)
	ErrUserNotFound        = errors.New(MsgUserNotFound)
	ErrPriorityOutOfRange  = errors.New(MsgPriorityOutOfRange)
	ErrInvalidStatus       = errors.New(MsgInvalidStatus)
	ErrJobAlreadyExists    = errors.New(MsgJobAlreadyExists)
)

func NewJob(title, description string, priority JobPriority, userID string) Job {
	return Job{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		Priority:    priority,
		Status:      JobStatusPending,
		UserID:      userID,
		CreatedAt:   time.Now().UTC().Truncate(time.Microsecond),
		UpdatedAt:   time.Now().UTC().Truncate(time.Microsecond),
	}
}

func (s JobStatus) IsValid() bool {
	switch s {
	case JobStatusPending, JobStatusRunning, JobStatusDone, JobStatusFailed:
		return true
	default:
		return false
	}
}

var validTransitions = map[JobStatus][]JobStatus{
	JobStatusPending: {JobStatusRunning, JobStatusFailed},
	JobStatusRunning: {JobStatusDone, JobStatusFailed},
	JobStatusFailed:  {JobStatusPending},
	JobStatusDone:    {}, // terminal
}

func (j *Job) CanTransitionTo(newStatus JobStatus) bool {
	for _, s := range validTransitions[j.Status] {
		if s == newStatus {
			return true
		}
	}
	return false
}
