package main

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

type JobService struct {
	repo Repository
}

func newJobService(repo Repository) *JobService {
	js := &JobService{repo}
	return js
}

func (js *JobService) GetByID(id string) (*Job, error) {
	job := js.repo.GetByID(id)
	if job == nil {
		return nil, ErrJobNotFound
	}
	return job, nil
}

func (js *JobService) Create(jobReq CreateJobRequest) *Job {
	nj := NewJob(jobReq.Title, jobReq.Description, JobPriority(jobReq.Priority), jobReq.UserID)
	js.repo.Create(&nj)
	return &nj
}

func (js *JobService) Update(jobID string, jobReq UpdateJobRequest) (*Job, error) {
	job := js.repo.GetByID(jobID)

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

	updatedJob := js.repo.Update(job)
	return updatedJob, nil
}

func (js *JobService) Delete(jobID string) error {
	return js.repo.Delete(jobID)
}
