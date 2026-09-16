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
	fmt.Println(row)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		fmt.Printf("error type: %T value %#v\n", err, err)
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
	return nil
}

func (pr *PostgresRepository) GetByID(userID string) (User, error) {
	return User{}, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
