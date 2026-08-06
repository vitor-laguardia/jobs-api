package main

import "time"

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
