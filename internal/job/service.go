package job

import (
	"errors"
	"time"
)

const (
	MsgInvalidStatusTransition = "invalid job status transition"
	MsgJobNotFound             = "job not found"
)

var (
	ErrJobNotFound         = errors.New(MsgJobNotFound)
	ErrJobStatusTransition = errors.New(MsgInvalidStatusTransition)
)

type Repository interface {
	GetByID(id string) *Job
	Create(job *Job)
	Update(job *Job) *Job
	Delete(jobID string) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	s := &Service{repo}
	return s
}

func (s *Service) GetByID(id string) (*Job, error) {
	job := s.repo.GetByID(id)
	if job == nil {
		return nil, ErrJobNotFound
	}
	return job, nil
}

func (s *Service) Create(jobReq CreateJobRequest) *Job {
	nj := NewJob(jobReq.Title, jobReq.Description, JobPriority(jobReq.Priority), jobReq.UserID)
	s.repo.Create(&nj)
	return &nj
}

func (s *Service) Update(jobID string, jobReq UpdateJobRequest) (*Job, error) {
	job := s.repo.GetByID(jobID)

	if job == nil {
		return nil, ErrJobNotFound
	}

	if jobReq.Title != "" {
		job.Title = jobReq.Title
	}
	if jobReq.Description != "" {
		job.Description = jobReq.Description
	}
	if jobReq.Priority != 0 {
		job.Priority = JobPriority(jobReq.Priority)
	}
	if jobReq.Status != "" {
		if !job.CanTransitionTo(JobStatus(jobReq.Status)) {
			return nil, ErrJobStatusTransition
		}
		job.Status = JobStatus(jobReq.Status)
	}
	job.UpdatedAt = time.Now()

	updatedJob := s.repo.Update(job)
	return updatedJob, nil
}

func (s *Service) Delete(jobID string) error {
	return s.repo.Delete(jobID)
}
