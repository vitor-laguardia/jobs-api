package user

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

func (pr *PostgresRepository) Create(user User) (User, error) {
	query := `
INSERT INTO users(id, name, email, created_at) 
VALUES($1, $2, $3, $4)
RETURNING id, name, email, created_at
`
	row := pr.DB.QueryRow(query, user.ID, user.Name, user.Email, user.CreatedAt)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return User{}, ErrDuplicateEmail
		}

		return User{}, fmt.Errorf("userRepository.Create: %w", err)
	}
	return user, nil
}

func (pr *PostgresRepository) Update(user User) (User, error) {
	return User{}, nil
}

func (pr *PostgresRepository) Delete(userID string) error {
	query := `
DELETE FROM users
WHERE id = $1
`
	result, err := pr.DB.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("userRepository.Delete: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("userRepository.Delete - RowsAffected: %w", err)
	}

	fmt.Printf("\n ROWS AFFECTED %d", rowsAffected)

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (pr *PostgresRepository) GetByID(userID string) (User, error) {
	query := `
SELECT id, name, email, created_at
FROM users
WHERE id = $1
`
	row := pr.DB.QueryRow(query, userID)
	user := User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)

	if err == sql.ErrNoRows {
		return User{}, ErrNotFound
	} else if err != nil {
		return User{}, fmt.Errorf("userRepository.GetByID: %w", err)
	}

	return user, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
