package userrepo

import (
	"errors"
	"fmt"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrDuplicateEmail  = errors.New("email address already registered")
	ErrInvalidUsername = errors.New("username cannot be empty")
)

type User struct {
	ID       int
	Username string
	Email    string
}

type UserRepository struct {
	users map[int]User
}

func NewRepository() *UserRepository {
	return &UserRepository{
		users: make(map[int]User),
	}
}

func (r *UserRepository) CreateUser(id int, username, email string) error {
	if username == "" {
		return ErrInvalidUsername
	}
	for _, u := range r.users {
		if u.Email == email {
			return ErrDuplicateEmail
		}
	}
	r.users[id] = User{ID: id, Username: username, Email: email}
	return nil
}

func (r *UserRepository) FindByID(id int) (*User, error) {
	u, exists := r.users[id]
	if !exists {
		return nil, fmt.Errorf("user lookup id %d: %w", id, ErrUserNotFound)
	}
	return &u, nil
}
