package user

import (
	"errors"
	"time"
)

const (
	MsgNotFound = "user not found"
)

var (
	ErrNotFound = errors.New(MsgNotFound)
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
