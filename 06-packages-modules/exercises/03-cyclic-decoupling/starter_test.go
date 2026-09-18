package decoupling_test

import (
	"errors"
	"testing"

	"github.com/ashraftryfie/golang-experiment/06-packages-modules/exercises/03-cyclic-decoupling/models"
	"github.com/ashraftryfie/golang-experiment/06-packages-modules/exercises/03-cyclic-decoupling/order"
	"github.com/ashraftryfie/golang-experiment/06-packages-modules/exercises/03-cyclic-decoupling/user"
)

func TestDecoupledServices(t *testing.T) {
	userService := user.NewService()
	userService.AddUser(models.User{ID: "usr_1", Name: "Ashraf"})

	// Inject userService as the UserChecker implementation into order.NewService
	orderService := order.NewService(userService)

	// Valid order
	err := orderService.PlaceOrder(models.Order{
		ID:     "ord_100",
		UserID: "usr_1",
		Amount: 49.99,
	})
	if err != nil {
		t.Fatalf("expected order placement to succeed, got %v", err)
	}

	// Invalid order for missing user
	err = orderService.PlaceOrder(models.Order{
		ID:     "ord_101",
		UserID: "non_existent_user",
		Amount: 19.99,
	})
	if !errors.Is(err, order.ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}
