package task

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestHandler() *Handler {
	repo := NewMemoryRepository()
	service := NewService(repo)
	return NewHandler(service)
}

func TestHandlerCreateTask(t *testing.T) {
	handler := newTestHandler()

	body := `{
		"title": "Learn Handler Test",
		"description": "Testing create task handler",
		"status": "todo",
		"category": "backend"
	}`

	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.CreateTask(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "Learn Handler Test") {
		t.Errorf("expected response body to contain task title, got %s", rr.Body.String())
	}
}

func TestHandlerCreateTaskInvalidBody(t *testing.T) {
	handler := newTestHandler()

	body := `{"title":`

	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.CreateTask(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "invalid request body") {
		t.Errorf("expected invalid request body error, got %s", rr.Body.String())
	}
}

func TestHandlerCreateTaskValidationError(t *testing.T) {
	handler := newTestHandler()

	body := `{
		"title": "",
		"status": "todo"
	}`

	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.CreateTask(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "title is required") {
		t.Errorf("expected title required error, got %s", rr.Body.String())
	}
}

func TestHandlerGetTasks(t *testing.T) {
	handler := newTestHandler()

	createBody := `{
		"title": "Learn Get Tasks",
		"description": "Testing get tasks handler",
		"status": "todo",
		"category": "backend"
	}`

	createReq := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRR := httptest.NewRecorder()

	handler.CreateTask(createRR, createReq)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rr := httptest.NewRecorder()

	handler.GetTasks(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "Learn Get Tasks") {
		t.Errorf("expected response body to contain task title, got %s", rr.Body.String())
	}
}

func TestHandlerGetTasksWithFilter(t *testing.T) {
	handler := newTestHandler()

	taskOne := `{
		"title": "Learn Go",
		"description": "Backend task",
		"status": "todo",
		"category": "backend"
	}`

	taskTwo := `{
		"title": "Buy Milk",
		"description": "Personal task",
		"status": "done",
		"category": "personal"
	}`

	req1 := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(taskOne))
	req1.Header.Set("Content-Type", "application/json")
	rr1 := httptest.NewRecorder()
	handler.CreateTask(rr1, req1)

	req2 := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(taskTwo))
	req2.Header.Set("Content-Type", "application/json")
	rr2 := httptest.NewRecorder()
	handler.CreateTask(rr2, req2)

	req := httptest.NewRequest(http.MethodGet, "/tasks?status=done", nil)
	rr := httptest.NewRecorder()

	handler.GetTasks(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "Buy Milk") {
		t.Errorf("expected response to contain Buy Milk, got %s", rr.Body.String())
	}

	if strings.Contains(rr.Body.String(), "Learn Go") {
		t.Errorf("expected response not to contain Learn Go, got %s", rr.Body.String())
	}
}

func TestHandlerGetTaskByID(t *testing.T) {
	handler := newTestHandler()

	body := `{
		"title": "Get By ID Task",
		"description": "Testing get by id",
		"status": "todo",
		"category": "backend"
	}`

	createReq := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(body))
	createReq.Header.Set("Content-Type", "application/json")
	createRR := httptest.NewRecorder()

	handler.CreateTask(createRR, createReq)

	req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	req.SetPathValue("id", "1")

	rr := httptest.NewRecorder()

	handler.GetTaskByID(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "Get By ID Task") {
		t.Errorf("expected response body to contain task title, got %s", rr.Body.String())
	}
}

func TestHandlerGetTaskByIDNotFound(t *testing.T) {
	handler := newTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/tasks/999", nil)
	req.SetPathValue("id", "999")

	rr := httptest.NewRecorder()

	handler.GetTaskByID(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "task not found") {
		t.Errorf("expected task not found error, got %s", rr.Body.String())
	}
}

func TestHandlerGetTaskByIDInvalidID(t *testing.T) {
	handler := newTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/tasks/abc", nil)
	req.SetPathValue("id", "abc")

	rr := httptest.NewRecorder()

	handler.GetTaskByID(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "invalid task id") {
		t.Errorf("expected invalid task id error, got %s", rr.Body.String())
	}
}

func TestHandlerUpdateTaskByID(t *testing.T) {
	handler := newTestHandler()

	createBody := `{
		"title": "Old Task",
		"description": "Old description",
		"status": "todo",
		"category": "backend"
	}`

	createReq := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRR := httptest.NewRecorder()

	handler.CreateTask(createRR, createReq)

	updateBody := `{
		"title": "Updated Task",
		"description": "Updated description",
		"status": "done",
		"category": "learning"
	}`

	req := httptest.NewRequest(http.MethodPut, "/tasks/1", strings.NewReader(updateBody))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "1")

	rr := httptest.NewRecorder()

	handler.UpdateTaskByID(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "Updated Task") {
		t.Errorf("expected response body to contain updated title, got %s", rr.Body.String())
	}

	if !strings.Contains(rr.Body.String(), "done") {
		t.Errorf("expected response body to contain updated status, got %s", rr.Body.String())
	}
}

func TestHandlerDeleteTaskByID(t *testing.T) {
	handler := newTestHandler()

	createBody := `{
		"title": "Task To Delete",
		"description": "Testing delete",
		"status": "todo",
		"category": "backend"
	}`

	createReq := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRR := httptest.NewRecorder()

	handler.CreateTask(createRR, createReq)

	req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	req.SetPathValue("id", "1")

	rr := httptest.NewRecorder()

	handler.RemoveTaskByID(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "task deleted successfully") {
		t.Errorf("expected success message, got %s", rr.Body.String())
	}
}

func TestHandlerDeleteTaskByIDNotFound(t *testing.T) {
	handler := newTestHandler()

	req := httptest.NewRequest(http.MethodDelete, "/tasks/999", nil)
	req.SetPathValue("id", "999")

	rr := httptest.NewRecorder()

	handler.RemoveTaskByID(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "task not found") {
		t.Errorf("expected task not found error, got %s", rr.Body.String())
	}
}
