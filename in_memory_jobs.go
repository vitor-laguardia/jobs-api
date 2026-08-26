package main

type JobRepository struct {
	jobs map[string]*Job
}

func (r *JobRepository) GetByID(id string) *Job {
	return r.jobs[id]
}

func (r *JobRepository) Create(job *Job) {
	r.jobs[job.ID] = job
}

func (r *JobRepository) Update(job *Job) *Job {
	return nil
}

func NewInMemoryJobs() *JobRepository {
	return &JobRepository{make(map[string]*Job)}
}
