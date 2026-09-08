package job

type JobRepository struct {
	jobs map[string]*Job
}

func (jr *JobRepository) GetByID(id string) *Job {
	return jr.jobs[id]
}

func (jr *JobRepository) Create(job *Job) {
	jr.jobs[job.ID] = job
}

func (jr *JobRepository) Update(job *Job) *Job {
	if _, exists := jr.jobs[job.ID]; !exists {
		return nil
	}
	jr.jobs[job.ID] = job
	return job
}

func (jr *JobRepository) Delete(jobID string) error {
	if _, exists := jr.jobs[jobID]; !exists {
		return ErrJobNotFound
	}
	delete(jr.jobs, jobID)
	return nil
}

func NewInMemoryJobs() *JobRepository {
	return &JobRepository{make(map[string]*Job)}
}
