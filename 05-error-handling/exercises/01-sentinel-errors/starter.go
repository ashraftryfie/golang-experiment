package userrepo

import (
	"errors"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrDuplicateEmail  = errors.New("email address already registered")
	ErrInvalidUsername = errors.New("username cannot be empty")
	ErrNotImplemented  = errors.New("TODO: implement")
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
	// TODO: Validate username and email uniqueness, store in map
	return ErrNotImplemented
}

func (r *UserRepository) FindByID(id int) (*User, error) {
	// TODO: Lookup user in map, wrap ErrUserNotFound if missing
	return nil, ErrNotImplemented
}
