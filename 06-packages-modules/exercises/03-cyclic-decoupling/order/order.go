package order

import (
	"errors"

	"github.com/ashraftryfie/golang-experiment/06-packages-modules/exercises/03-cyclic-decoupling/models"
)

var ErrUserNotFound = errors.New("user does not exist")

// UserChecker is defined here at the CONSUMER boundary!
type UserChecker interface {
	UserExists(userID string) bool
}

type OrderService struct {
	checker UserChecker
	orders  map[string]models.Order
}

func NewService(checker UserChecker) *OrderService {
	return &OrderService{
		checker: checker,
		orders:  make(map[string]models.Order),
	}
}

func (s *OrderService) PlaceOrder(o models.Order) error {
	if s.checker == nil || !s.checker.UserExists(o.UserID) {
		return ErrUserNotFound
	}
	s.orders[o.ID] = o
	return nil
}
