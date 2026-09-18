package migrationrunner

import (
	"context"
	"errors"
	"testing"
)

type mockMigrationStore struct {
	tableCreated bool
	applied      map[int]bool
	history      []int
	failOnVer    int
}

func newMockStore() *mockMigrationStore {
	return &mockMigrationStore{
		applied: make(map[int]bool),
		history: make([]int, 0),
	}
}

func (m *mockMigrationStore) EnsureSchemaTable(ctx context.Context) error {
	m.tableCreated = true
	return nil
}

func (m *mockMigrationStore) GetAppliedVersions(ctx context.Context) (map[int]bool, error) {
	copyMap := make(map[int]bool)
	for k, v := range m.applied {
		copyMap[k] = v
	}
	return copyMap, nil
}

func (m *mockMigrationStore) ApplyMigration(ctx context.Context, mig Migration) error {
	if mig.Version == m.failOnVer {
		return errors.New("simulated SQL syntax error in migration")
	}
	m.applied[mig.Version] = true
	m.history = append(m.history, mig.Version)
	return nil
}

func TestRunMigrationsLifecycle(t *testing.T) {
	ctx := context.Background()
	store := newMockStore()

	migrations := []Migration{
		{Version: 1, Description: "create users table", SQL: "CREATE TABLE users (...);"},
		{Version: 2, Description: "create orders table", SQL: "CREATE TABLE orders (...);"},
		{Version: 3, Description: "add index on users", SQL: "CREATE INDEX idx_users (...);"},
	}

	// 1. First run: all 3 should be applied
	count, err := RunMigrations(ctx, store, migrations)
	if err != nil {
		t.Fatalf("unexpected error running migrations: %v", err)
	}

	if count != 3 {
		t.Errorf("expected 3 migrations applied, got %d", count)
	}
	if !store.tableCreated {
		t.Errorf("expected schema table to be created")
	}

	// 2. Second run: idempotency - 0 should be applied
	count2, err := RunMigrations(ctx, store, migrations)
	if err != nil {
		t.Fatalf("unexpected error running migrations second time: %v", err)
	}
	if count2 != 0 {
		t.Errorf("expected 0 migrations applied on re-run, got %d", count2)
	}

	// 3. New migration added (version 4)
	extendedMigrations := append(migrations, Migration{
		Version:     4,
		Description: "add payments table",
		SQL:         "CREATE TABLE payments (...);",
	})
	count3, err := RunMigrations(ctx, store, extendedMigrations)
	if err != nil {
		t.Fatalf("unexpected error applying v4: %v", err)
	}
	if count3 != 1 {
		t.Errorf("expected 1 migration applied for v4, got %d", count3)
	}
}

func TestRunMigrationsErrors(t *testing.T) {
	ctx := context.Background()

	// Out of order
	badOrder := []Migration{
		{Version: 2, Description: "v2", SQL: ""},
		{Version: 1, Description: "v1", SQL: ""},
	}
	_, err := RunMigrations(ctx, newMockStore(), badOrder)
	if err == nil {
		t.Errorf("expected error for unsorted migrations, got nil")
	}

	// Duplicate version
	dups := []Migration{
		{Version: 1, Description: "v1", SQL: ""},
		{Version: 1, Description: "v1 dupe", SQL: ""},
	}
	_, err = RunMigrations(ctx, newMockStore(), dups)
	if err == nil {
		t.Errorf("expected error for duplicate migration versions, got nil")
	}

	// Failure aborts remaining
	storeFail := newMockStore()
	storeFail.failOnVer = 2
	all := []Migration{
		{Version: 1, Description: "v1", SQL: ""},
		{Version: 2, Description: "v2", SQL: ""},
		{Version: 3, Description: "v3", SQL: ""},
	}
	_, err = RunMigrations(ctx, storeFail, all)
	if err == nil {
		t.Errorf("expected error on failing migration v2, got nil")
	}
	if storeFail.applied[3] {
		t.Errorf("migration v3 should not have been attempted after v2 failed")
	}
}
