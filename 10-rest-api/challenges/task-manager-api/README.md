# Challenge: Production-Grade Task Management REST API

## 🎯 Objective
Build a RESTful task management API supporting full CRUD semantics, query filtering, pagination, payload validation, and thread-safe persistence using Go 1.22+ `http.ServeMux`.

## 📋 Requirements
1. **Domain Model**:
   ```go
   type Task struct {
       ID        string    `json:"id"`
       Title     string    `json:"title"`
       Status    string    `json:"status"`   // "pending" or "completed"
       Priority  int       `json:"priority"` // 1 (lowest) to 5 (highest)
       CreatedAt time.Time `json:"created_at"`
   }
   ```

2. **Endpoints**:
   - `POST /tasks`: Validates `title` (non-empty) and `priority` (1-5). Returns `201 Created` with `Location: /tasks/{id}`.
   - `GET /tasks`: Supports `?status=pending|completed`, `?page=1`, and `?limit=10`. Returns:
     ```json
     {
       "data": [ ... ],
       "pagination": { "page": 1, "limit": 10, "total": 42 }
     }
     ```
   - `GET /tasks/{id}`: Returns `200 OK` or `404 Not Found`.
   - `PUT /tasks/{id}`: Replaces task; returns `200 OK` or `404 Not Found`.
   - `DELETE /tasks/{id}`: Removes task; returns `204 No Content` or `404 Not Found`.

3. **Concurrency Safety**:
   All state access must be guarded by `sync.RWMutex`.

## 🧪 Verification
```powershell
go test -v ./10-rest-api/challenges/task-manager-api/...
```
