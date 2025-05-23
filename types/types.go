package types

import (
	"time"
)

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(ID int) (*User, error)
	CreateUser(User) error
}

type User struct {
	ID        int
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreateAt  time.Time `json:"createAt"`
}

type RegisterUserPayload struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}
