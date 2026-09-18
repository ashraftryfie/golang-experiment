# Challenge: Transactional Database Migration Runner

## 🎯 Objective
Build a transactional database migration engine that validates schema evolution scripts, tracks applied versions in a metadata table, executes DDL scripts in order, and aborts safely on error.

## 📋 Requirements
1. **Migration Definition**:
   ```go
   type Migration struct {
       Version     int
       Description string
       SQL         string
   }
   ```

2. **MigrationStore Abstraction**:
   ```go
   type MigrationStore interface {
       EnsureSchemaTable(ctx context.Context) error
       GetAppliedVersions(ctx context.Context) (map[int]bool, error)
       ApplyMigration(ctx context.Context, m Migration) error
   }
   ```

3. **Execution Logic**:
   - `RunMigrations(ctx context.Context, store MigrationStore, migrations []Migration) (int, error)`
   - Validates that `migrations` versions are strictly increasing positive integers ($v_1 < v_2 < v_3$).
   - Calls `EnsureSchemaTable` to create `schema_migrations` tracking table if missing.
   - Fetches already applied version set.
   - Applies pending migrations sequentially.
   - If any migration fails, aborts immediately and propagates the error.
   - Returns the count of newly executed migrations.

## 🧪 Verification
```powershell
go test -v ./11-postgresql/challenges/migration-runner/...
```
