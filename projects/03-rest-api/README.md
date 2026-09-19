# Project 3: `task-service` (Task Service REST API)

A clean, production-styled RESTful task management service featuring request validation, clean routing with standard library `http.ServeMux` (Go 1.22+), structured error responses, and comprehensive HTTP integration tests.

## Endpoints
- `POST /tasks`: Create a new task (`title`, `description`, `status`)
- `GET /tasks`: List all tasks (supports query filtering)
- `GET /tasks/{id}`: Get single task by ID
- `PUT /tasks/{id}`: Update task status/content
- `DELETE /tasks/{id}`: Remove task

## Usage
```powershell
go test -v ./projects/03-rest-api/...
```
