package job

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

type PostgresRepository struct {
	DB *sql.DB
}

func NewPostgresRepository(DB *sql.DB) *PostgresRepository {
	return &PostgresRepository{DB: DB}
}

func (pr *PostgresRepository) Create(job *Job) (int, error) {
	query := `
INSERT INTO jobs(id, user_id, title, description, status, priority, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`
	result, err := pr.DB.Exec(query, job.ID, job.UserID, job.Title, job.Description, job.Status, job.Priority, job.CreatedAt, job.UpdatedAt)

	if err != nil {
		switch {
		case isFKViolation(err):
			return 0, ErrUserNotFound
		case isJobUniqueConstraintViolation(err):
			return 0, ErrJobAlreadyExists
		default:
			return 0, fmt.Errorf("insert job: %w", err)
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("get rows affected: %w", err)
	}
	return int(rowsAffected), nil
}

func (pr *PostgresRepository) GetByID(id string) (*Job, error) {
	query := `
SELECT id, title, description, status, priority, user_id, created_at, updated_at
FROM jobs
WHERE id = $1
`
	row := pr.DB.QueryRow(query, id)
	job := Job{}
	err := row.Scan(&job.ID, &job.Title, &job.Description, &job.Status, &job.Priority, &job.UserID, &job.CreatedAt, &job.UpdatedAt)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrJobNotFound
		default:
			return nil, fmt.Errorf("retrieve job: %w", err)
		}
	}

	job.CreatedAt = job.CreatedAt.UTC()
	job.UpdatedAt = job.UpdatedAt.UTC()

	return &job, nil
}

func (pr *PostgresRepository) Update(job *Job) (*Job, error) {
	query := `
UPDATE jobs
SET title = $1,
    description = $2,
    status = $3,
    priority = $4,
    updated_at = now()
WHERE id = $5
RETURNING id, title, description, status, priority, user_id, created_at, updated_at
`
	newJob := new(Job)
	row := pr.DB.QueryRow(query, job.Title, job.Description, job.Status, job.Priority, job.ID)
	err := row.Scan(&newJob.ID, &newJob.Title, &newJob.Description, &newJob.Status, &newJob.Priority, &newJob.UserID, &newJob.CreatedAt, &newJob.UpdatedAt)

	if err != nil {
		switch {
		case isPriorityConstraintViolation(err):
			return nil, ErrPriorityOutOfRange
		case isStatusConstraintViolation(err):
			return nil, ErrInvalidStatus
		default:
			return nil, fmt.Errorf("update job: %w", err)
		}
	}

	newJob.CreatedAt = newJob.CreatedAt.UTC()
	newJob.UpdatedAt = newJob.UpdatedAt.UTC()

	return newJob, nil
}

func (pr *PostgresRepository) Delete(jobID string) error {
	query := `
DELETE from jobs
WHERE id = $1
`
	result, err := pr.DB.Exec(query, jobID)
	if err != nil {
		return fmt.Errorf("exclude job: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrJobNotFound
	}

	return nil
}

func isFKViolation(err error) bool {
	var pgError *pgconn.PgError
	return errors.As(err, &pgError) && pgError.Code == "23503"
}
