package job

import (
	"errors"
	"time"

	"github.com/vitor-laguardia/jobs-api/internal/user"
)

type JobRepository struct {
	jobs     map[string]*Job
	userRepo user.Repository
}

func NewInMemoryRepository(userRepo user.Repository) *JobRepository {
	return &JobRepository{
		jobs:     make(map[string]*Job),
		userRepo: userRepo,
	}
}

func (jr *JobRepository) Create(job *Job) (int, error) {
	_, err := jr.userRepo.GetByID(job.UserID)
	if errors.Is(err, user.ErrNotFound) {
		return 0, ErrUserNotFound
	}

	if _, exists := jr.jobs[job.ID]; exists {
		return 0, ErrJobAlreadyExists
	}

	jobCopy := *job
	jr.jobs[job.ID] = &jobCopy
	return 1, nil
}

func (jr *JobRepository) GetByID(id string) (*Job, error) {
	job, ok := jr.jobs[id]
	if !ok {
		return nil, ErrJobNotFound
	}

	jobCopy := *job
	return &jobCopy, nil
}

func (jr *JobRepository) Update(job *Job) (*Job, error) {
	if _, exists := jr.jobs[job.ID]; !exists {
		return nil, ErrJobNotFound
	}

	if !job.Priority.IsValid() {
		return nil, ErrPriorityOutOfRange
	}

	if !job.Status.IsValid() {
		return nil, ErrInvalidStatus
	}

	updated := *job
	updated.UpdatedAt = time.Now().UTC()
	jr.jobs[job.ID] = &updated

	result := updated
	return &result, nil
}

func (jr *JobRepository) Delete(jobID string) error {
	if _, exists := jr.jobs[jobID]; !exists {
		return ErrJobNotFound
	}

	delete(jr.jobs, jobID)
	return nil
}
