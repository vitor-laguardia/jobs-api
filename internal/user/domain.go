package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	MsgNotFound       = "user not found"
	MsgDuplicateEmail = "email already in use"
)

var (
	ErrNotFound       = errors.New(MsgNotFound)
	ErrDuplicateEmail = errors.New(MsgDuplicateEmail)
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func NewUser(name, email string) User {
	return User{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
	}
}
