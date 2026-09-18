package user

import (
	"github.com/ashraftryfie/golang-experiment/06-packages-modules/exercises/03-cyclic-decoupling/models"
)

type UserService struct {
	users map[string]models.User
}

func NewService() *UserService {
	return &UserService{
		users: make(map[string]models.User),
	}
}

func (s *UserService) AddUser(u models.User) {
	s.users[u.ID] = u
}

// UserExists satisfies the OrderService's UserChecker interface
func (s *UserService) UserExists(id string) bool {
	_, exists := s.users[id]
	return exists
}
