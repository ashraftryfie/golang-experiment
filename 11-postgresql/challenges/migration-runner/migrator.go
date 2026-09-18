package migrationrunner

import (
	"context"
	"errors"
	"fmt"
	"sort"
)

type Migration struct {
	Version     int
	Description string
	SQL         string
}

type MigrationStore interface {
	EnsureSchemaTable(ctx context.Context) error
	GetAppliedVersions(ctx context.Context) (map[int]bool, error)
	ApplyMigration(ctx context.Context, m Migration) error
}

// RunMigrations validates and executes pending database migrations in strict ascending order.
func RunMigrations(ctx context.Context, store MigrationStore, migrations []Migration) (int, error) {
	if err := validateMigrations(migrations); err != nil {
		return 0, fmt.Errorf("migration validation failed: %w", err)
	}

	if err := store.EnsureSchemaTable(ctx); err != nil {
		return 0, fmt.Errorf("ensure schema table failed: %w", err)
	}

	applied, err := store.GetAppliedVersions(ctx)
	if err != nil {
		return 0, fmt.Errorf("get applied versions failed: %w", err)
	}

	count := 0
	for _, m := range migrations {
		if applied[m.Version] {
			continue
		}

		if err := store.ApplyMigration(ctx, m); err != nil {
			return count, fmt.Errorf("apply migration v%d (%s) failed: %w", m.Version, m.Description, err)
		}
		count++
	}

	return count, nil
}

func validateMigrations(migrations []Migration) error {
	if len(migrations) == 0 {
		return nil
	}

	seen := make(map[int]bool)
	versions := make([]int, len(migrations))

	for i, m := range migrations {
		if m.Version <= 0 {
			return fmt.Errorf("version must be positive, got %d", m.Version)
		}
		if seen[m.Version] {
			return fmt.Errorf("duplicate migration version %d detected", m.Version)
		}
		seen[m.Version] = true
		versions[i] = m.Version
	}

	if !sort.IntsAreSorted(versions) {
		return errors.New("migrations must be strictly sorted by version ascending")
	}

	return nil
}
