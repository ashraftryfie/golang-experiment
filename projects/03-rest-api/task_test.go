package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTaskHandler_E2E(t *testing.T) {
	store := NewInMemoryTaskStore()
	handler := NewTaskHandler(store)
	server := httptest.NewServer(handler)
	defer server.Close()

	client := server.Client()

	// 1. POST /tasks
	postBody := `{"title":"Deploy to K8s","description":"Deploy service with Helm"}`
	resp, err := client.Post(server.URL+"/tasks", "application/json", bytes.NewBufferString(postBody))
	if err != nil {
		t.Fatalf("POST /tasks failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", resp.StatusCode)
	}

	var created Task
	_ = json.NewDecoder(resp.Body).Decode(&created)
	resp.Body.Close()

	if created.Title != "Deploy to K8s" || created.Status != StatusTodo {
		t.Fatalf("unexpected task state: %+v", created)
	}

	// 2. GET /tasks/{id}
	getResp, err := client.Get(server.URL + "/tasks/" + created.ID)
	if err != nil || getResp.StatusCode != http.StatusOK {
		t.Fatalf("GET /tasks/{id} failed: %v, status: %d", err, getResp.StatusCode)
	}
	getResp.Body.Close()

	// 3. PUT /tasks/{id}
	putBody := `{"title":"Deploy to K8s","description":"Updated description","status":"DONE"}`
	req, _ := http.NewRequest(http.MethodPut, server.URL+"/tasks/"+created.ID, bytes.NewBufferString(putBody))
	req.Header.Set("Content-Type", "application/json")
	putResp, err := client.Do(req)
	if err != nil || putResp.StatusCode != http.StatusOK {
		t.Fatalf("PUT /tasks/{id} failed: %v, status: %d", err, putResp.StatusCode)
	}
	var updated Task
	_ = json.NewDecoder(putResp.Body).Decode(&updated)
	putResp.Body.Close()

	if updated.Status != StatusDone {
		t.Fatalf("expected status DONE, got %s", updated.Status)
	}

	// 4. DELETE /tasks/{id}
	delReq, _ := http.NewRequest(http.MethodDelete, server.URL+"/tasks/"+created.ID, nil)
	delResp, err := client.Do(delReq)
	if err != nil || delResp.StatusCode != http.StatusNoContent {
		t.Fatalf("DELETE /tasks/{id} failed: %v, status: %d", err, delResp.StatusCode)
	}
	delResp.Body.Close()

	// 5. Verify deleted
	getDeleted, _ := client.Get(server.URL + "/tasks/" + created.ID)
	if getDeleted.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 for deleted task, got %d", getDeleted.StatusCode)
	}
	getDeleted.Body.Close()
}
