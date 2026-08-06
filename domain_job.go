package main

import (
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

func NewJob(title, description string, priority JobPriority, userID string) Job {
	return Job{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		Priority:    priority,
		Status:      JobStatusPending,
		UserID:      userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
