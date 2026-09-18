package main

import (
	"fmt"
	"log"

	"github.com/ashraftryfie/golang-experiment/06-packages-modules/challenges/modular-monolith/internal/domain"
	"github.com/ashraftryfie/golang-experiment/06-packages-modules/challenges/modular-monolith/internal/store"
	"github.com/ashraftryfie/golang-experiment/06-packages-modules/challenges/modular-monolith/pkg/validator"
)

type App struct {
	store store.Store
}

func NewApp(st store.Store) *App {
	return &App{store: st}
}

func (a *App) CreateProduct(id, name string, price float64) error {
	if err := validator.ValidateProduct(id, name, price); err != nil {
		return fmt.Errorf("invalid product: %w", err)
	}
	return a.store.Save(domain.Product{ID: id, Name: name, Price: price})
}

func main() {
	st := store.NewMemoryStore()
	app := NewApp(st)

	if err := app.CreateProduct("p100", "Mechanical Keyboard", 129.99); err != nil {
		log.Fatalf("failed: %v", err)
	}

	p, err := st.Get("p100")
	if err != nil {
		log.Fatalf("failed: %v", err)
	}

	fmt.Printf("Created Product: %s ($%.2f)\n", p.Name, p.Price)
}
