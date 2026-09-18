package taskmanager

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTaskAPIWorkflow(t *testing.T) {
	store := NewTaskStore()
	router := NewTaskRouter(store)

	// 1. Create Task (POST /tasks) -> 201 Created
	createBody := `{"title":"Ship Stage 10 REST API","priority":5}`
	reqCreate := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(createBody))
	recCreate := httptest.NewRecorder()
	router.ServeHTTP(recCreate, reqCreate)

	if recCreate.Code != http.StatusCreated {
		t.Fatalf("create failed: status %d, body %s", recCreate.Code, recCreate.Body.String())
	}

	var created Task
	if err := json.Unmarshal(recCreate.Body.Bytes(), &created); err != nil {
		t.Fatalf("failed to decode created task: %v", err)
	}

	if created.ID != "task-1" || created.Status != "pending" || created.Priority != 5 {
		t.Errorf("created task fields unexpected: %+v", created)
	}

	// 2. Query with pagination and status filter (GET /tasks?status=pending&page=1&limit=5)
	reqQuery := httptest.NewRequest(http.MethodGet, "/tasks?status=pending&page=1&limit=5", nil)
	recQuery := httptest.NewRecorder()
	router.ServeHTTP(recQuery, reqQuery)

	if recQuery.Code != http.StatusOK {
		t.Fatalf("query failed: status %d", recQuery.Code)
	}

	var pageResp PaginatedTasksResponse
	if err := json.Unmarshal(recQuery.Body.Bytes(), &pageResp); err != nil {
		t.Fatalf("failed to decode paginated tasks: %v", err)
	}

	if pageResp.Pagination.Total != 1 || len(pageResp.Data) != 1 {
		t.Errorf("expected 1 task in query result, got %+v", pageResp)
	}

	// 3. Update Task to completed (PUT /tasks/task-1)
	updateBody := `{"title":"Ship Stage 10 REST API","priority":5,"status":"completed"}`
	reqUpdate := httptest.NewRequest(http.MethodPut, "/tasks/task-1", bytes.NewBufferString(updateBody))
	recUpdate := httptest.NewRecorder()
	router.ServeHTTP(recUpdate, reqUpdate)

	if recUpdate.Code != http.StatusOK {
		t.Fatalf("update failed: status %d", recUpdate.Code)
	}

	// 4. Query with status=pending should now return 0 items
	reqPending := httptest.NewRequest(http.MethodGet, "/tasks?status=pending", nil)
	recPending := httptest.NewRecorder()
	router.ServeHTTP(recPending, reqPending)

	var pendingResp PaginatedTasksResponse
	_ = json.Unmarshal(recPending.Body.Bytes(), &pendingResp)
	if pendingResp.Pagination.Total != 0 || len(pendingResp.Data) != 0 {
		t.Errorf("expected 0 pending tasks after update, got %d", pendingResp.Pagination.Total)
	}

	// 5. Delete Task (DELETE /tasks/task-1) -> 204 No Content
	reqDelete := httptest.NewRequest(http.MethodDelete, "/tasks/task-1", nil)
	recDelete := httptest.NewRecorder()
	router.ServeHTTP(recDelete, reqDelete)

	if recDelete.Code != http.StatusNoContent {
		t.Fatalf("delete failed: status %d", recDelete.Code)
	}

	// 6. Verify 404 after delete
	reqGetMissing := httptest.NewRequest(http.MethodGet, "/tasks/task-1", nil)
	recGetMissing := httptest.NewRecorder()
	router.ServeHTTP(recGetMissing, reqGetMissing)

	if recGetMissing.Code != http.StatusNotFound {
		t.Errorf("expected 404 after delete, got %d", recGetMissing.Code)
	}
}

func TestValidationErrors(t *testing.T) {
	store := NewTaskStore()
	router := NewTaskRouter(store)

	// Invalid priority (0)
	badPriority := `{"title":"Valid Title","priority":0}`
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(badPriority))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422 for priority 0, got %d", rec.Code)
	}

	// Empty title
	emptyTitle := `{"title":"","priority":3}`
	req2 := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(emptyTitle))
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422 for empty title, got %d", rec2.Code)
	}
}
