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
		if isFKViolation(err) {
			return 0, ErrUserNotFound
		}

		return 0, fmt.Errorf("jobRepository.Create: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	return int(rowsAffected), nil
}

func (pr *PostgresRepository) GetByID(id string) (*Job, error) {
	return nil, nil
}

func (pr *PostgresRepository) Update(job *Job) (*Job, error) {
	return nil, nil
}

func (pr *PostgresRepository) Delete(jobID string) error {
	return nil
}

func isFKViolation(err error) bool {
	var pgError *pgconn.PgError
	return errors.As(err, &pgError) && pgError.Code == "23503"
}
