# Exercise 03: Docker Compose Configuration Validator

## Objective
Implement a static analysis validator for Docker Compose files to prevent unencrypted secrets from leaking into Git repositories and enforce dependency healthchecks.

## Requirements
1. **Required Services**: Ensure all specified `requiredServices` exist under the `services` top-level mapping (`MISSING_REQUIRED_SERVICE`).
2. **Plaintext Password Detection**: Flag any environment variable containing `PASSWORD` or `SECRET` whose value does NOT use variable interpolation `${...}` (`HARDCODED_SECRET_DETECTED`).
3. **Database Healthcheck**: Any database service (e.g. containing `postgres`, `mysql`, `mongo`, `redis`) must define a `healthcheck` block (`DATABASE_HEALTHCHECK_REQUIRED`).

## Starter & Tests
- Starter: `starter.go`
- Run starter tests: `go test -v ./18-docker/exercises/03-compose-config-validator`
