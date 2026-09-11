package job

type JobRepository struct {
	jobs map[string]*Job
}

func (jr *JobRepository) GetByID(id string) (*Job, error) {
	job, ok := jr.jobs[id]

	if !ok {
		return nil, ErrJobNotFound
	}

	jobCopy := *job
	return &jobCopy, nil
}

func (jr *JobRepository) Create(job *Job) error {
	jr.jobs[job.ID] = job
	return nil
}

func (jr *JobRepository) Update(job *Job) (*Job, error) {
	jr.jobs[job.ID] = job
	return job, nil
}

func (jr *JobRepository) Delete(jobID string) error {
	if _, exists := jr.jobs[jobID]; !exists {
		return ErrJobNotFound
	}
	delete(jr.jobs, jobID)
	return nil
}

func NewInMemoryRepository() *JobRepository {
	return &JobRepository{make(map[string]*Job)}
}
