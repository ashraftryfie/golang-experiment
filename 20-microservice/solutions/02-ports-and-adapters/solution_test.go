package portsadapters

import (
	"context"
	"errors"
	"sync"
	"testing"
)

func TestProduct_Invariants(t *testing.T) {
	p, err := NewProduct("p1", "SKU-TEST", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := p.DeductStock(3); err != nil {
		t.Fatalf("unexpected deduction error: %v", err)
	}
	if p.Stock() != 7 {
		t.Fatalf("expected stock 7, got %d", p.Stock())
	}

	if err := p.DeductStock(10); !errors.Is(err, ErrInsufficientStock) {
		t.Fatalf("expected ErrInsufficientStock, got %v", err)
	}

	if err := p.AddStock(5); err != nil {
		t.Fatalf("unexpected add error: %v", err)
	}
	if p.Stock() != 12 {
		t.Fatalf("expected stock 12, got %d", p.Stock())
	}
}

func TestInMemoryProductRepo_ThreadSafetyAndIsolation(t *testing.T) {
	ctx := context.Background()
	repo := NewInMemoryProductRepo()

	prod, _ := NewProduct("p1", "SKU-LAPTOP", 100)
	if err := repo.Save(ctx, prod); err != nil {
		t.Fatalf("failed to save product: %v", err)
	}

	// Verify defensive copy on retrieve
	fetched, err := repo.GetBySKU(ctx, "SKU-LAPTOP")
	if err != nil {
		t.Fatalf("failed to fetch product: %v", err)
	}
	_ = fetched.DeductStock(50) // modify local copy

	// Re-fetch to ensure repository state was not mutated externally
	unmodified, err := repo.GetBySKU(ctx, "SKU-LAPTOP")
	if err != nil {
		t.Fatalf("failed to re-fetch product: %v", err)
	}
	if unmodified.Stock() != 100 {
		t.Fatalf("repository state was mutated directly! expected 100, got %d", unmodified.Stock())
	}

	// Concurrent read/write stress test
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			_, _ = repo.GetBySKU(ctx, "SKU-LAPTOP")
		}()
		go func(val int) {
			defer wg.Done()
			p, _ := NewProduct("p1", "SKU-LAPTOP", val)
			_ = repo.Save(ctx, p)
		}(i)
	}
	wg.Wait()
}

func TestMockAlertPublisher(t *testing.T) {
	ctx := context.Background()
	pub := NewMockAlertPublisher()

	if err := pub.PublishLowStock(ctx, "SKU-A", 2); err != nil {
		t.Fatalf("publish failed: %v", err)
	}
	if err := pub.PublishLowStock(ctx, "SKU-B", 0); err != nil {
		t.Fatalf("publish failed: %v", err)
	}

	alerts := pub.RecordedAlerts()
	if len(alerts) != 2 || alerts["SKU-A"] != 2 || alerts["SKU-B"] != 0 {
		t.Fatalf("unexpected alerts recorded: %v", alerts)
	}
}
