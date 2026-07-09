package task

import (
	"errors"
	"testing"
)

func newTestService() *Service {
	repo := NewMemoryRepository()
	return NewService(repo)
}

func TestServiceCreateTask(t *testing.T) {
	service := newTestService()

	req := CreateTaskRequest{
		Title:       "Learn Go API",
		Description: "Build a REST API with Go",
		Status:      "todo",
		Category:    "backend",
	}

	task, err := service.CreateTask(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if task.ID != 1 {
		t.Errorf("expected task ID 1, got %d", task.ID)
	}

	if task.Title != "Learn Go API" {
		t.Errorf("expected title %q, got %q", "Learn Go API", task.Title)
	}

	if task.Status != "todo" {
		t.Errorf("expected status %q, got %q", "todo", task.Status)
	}

	if task.CreatedAt == "" {
		t.Error("expected CreatedAt to be set")
	}

	if task.UpdatedAt == "" {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestServiceCreateTaskWithDefaultStatus(t *testing.T) {
	service := newTestService()

	req := CreateTaskRequest{
		Title: "Learn Testing",
	}

	task, err := service.CreateTask(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if task.Status != "todo" {
		t.Errorf("expected default status %q, got %q", "todo", task.Status)
	}
}

func TestServiceCreateTaskValidationError(t *testing.T) {
	service := newTestService()

	tests := []struct {
		name      string
		req       CreateTaskRequest
		wantError error
	}{
		{
			name: "empty title",
			req: CreateTaskRequest{
				Title: "",
			},
			wantError: ErrTitleRequired,
		},
		{
			name: "short title",
			req: CreateTaskRequest{
				Title: "Go",
			},
			wantError: ErrTitleTooShort,
		},
		{
			name: "invalid status",
			req: CreateTaskRequest{
				Title:  "Learn Go",
				Status: "finished",
			},
			wantError: ErrInvalidStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.CreateTask(tt.req)

			if !errors.Is(err, tt.wantError) {
				t.Errorf("expected error %v, got %v", tt.wantError, err)
			}
		})
	}
}

func TestServiceGetTaskByID(t *testing.T) {
	service := newTestService()

	createdTask, err := service.CreateTask(CreateTaskRequest{
		Title:  "Learn Go",
		Status: "todo",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	task, err := service.GetTaskByID(createdTask.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if task.ID != createdTask.ID {
		t.Errorf("expected task ID %d, got %d", createdTask.ID, task.ID)
	}
}

func TestServiceGetTaskByIDNotFound(t *testing.T) {
	service := newTestService()

	_, err := service.GetTaskByID(999)

	if !errors.Is(err, ErrTaskNotFound) {
		t.Errorf("expected error %v, got %v", ErrTaskNotFound, err)
	}
}

func TestServiceUpdateTask(t *testing.T) {
	service := newTestService()

	createdTask, err := service.CreateTask(CreateTaskRequest{
		Title:       "Learn Go",
		Description: "Old description",
		Status:      "todo",
		Category:    "backend",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updatedTask, err := service.UpdateTask(createdTask.ID, UpdateTaskRequest{
		Title:       "Learn Advanced Go",
		Description: "Updated description",
		Status:      "done",
		Category:    "learning",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if updatedTask.Title != "Learn Advanced Go" {
		t.Errorf("expected updated title, got %q", updatedTask.Title)
	}

	if updatedTask.Status != "done" {
		t.Errorf("expected status done, got %q", updatedTask.Status)
	}

	if updatedTask.Category != "learning" {
		t.Errorf("expected category learning, got %q", updatedTask.Category)
	}

	//if updatedTask.UpdatedAt == createdTask.UpdatedAt {
	//	t.Error("expected UpdatedAt to change")
	//}
}

func TestServiceUpdateTaskNotFound(t *testing.T) {
	service := newTestService()

	_, err := service.UpdateTask(999, UpdateTaskRequest{
		Title: "Updated Task",
	})

	if !errors.Is(err, ErrTaskNotFound) {
		t.Errorf("expected error %v, got %v", ErrTaskNotFound, err)
	}
}

func TestServiceDeleteTask(t *testing.T) {
	service := newTestService()

	createdTask, err := service.CreateTask(CreateTaskRequest{
		Title: "Task to delete",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = service.DeleteTask(createdTask.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = service.GetTaskByID(createdTask.ID)
	if !errors.Is(err, ErrTaskNotFound) {
		t.Errorf("expected task to be deleted, got error %v", err)
	}
}

func TestServiceDeleteTaskNotFound(t *testing.T) {
	service := newTestService()

	err := service.DeleteTask(999)

	if !errors.Is(err, ErrTaskNotFound) {
		t.Errorf("expected error %v, got %v", ErrTaskNotFound, err)
	}
}

func TestServiceGetTasksWithFilters(t *testing.T) {
	service := newTestService()

	_, err := service.CreateTask(CreateTaskRequest{
		Title:       "Learn Go",
		Description: "Build REST API",
		Status:      "todo",
		Category:    "backend",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = service.CreateTask(CreateTaskRequest{
		Title:       "Buy Milk",
		Description: "Personal task",
		Status:      "done",
		Category:    "personal",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = service.CreateTask(CreateTaskRequest{
		Title:       "Learn PostgreSQL",
		Description: "Database for Go API",
		Status:      "in_progress",
		Category:    "backend",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	t.Run("filter by status", func(t *testing.T) {
		tasks, err := service.GetTasks(TaskFilter{
			Status: "done",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(tasks) != 1 {
			t.Fatalf("expected 1 task, got %d", len(tasks))
		}

		if tasks[0].Status != "done" {
			t.Errorf("expected status done, got %q", tasks[0].Status)
		}
	})

	t.Run("filter by category", func(t *testing.T) {
		tasks, err := service.GetTasks(TaskFilter{
			Category: "backend",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(tasks) != 2 {
			t.Fatalf("expected 2 tasks, got %d", len(tasks))
		}
	})

	t.Run("search by title or description", func(t *testing.T) {
		tasks, err := service.GetTasks(TaskFilter{
			Search: "postgresql",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(tasks) != 1 {
			t.Fatalf("expected 1 task, got %d", len(tasks))
		}

		if tasks[0].Title != "Learn PostgreSQL" {
			t.Errorf("expected Learn PostgreSQL, got %q", tasks[0].Title)
		}
	})

	t.Run("combined filter", func(t *testing.T) {
		tasks, err := service.GetTasks(TaskFilter{
			Status:   "todo",
			Category: "backend",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(tasks) != 1 {
			t.Fatalf("expected 1 task, got %d", len(tasks))
		}

		if tasks[0].Title != "Learn Go" {
			t.Errorf("expected Learn Go, got %q", tasks[0].Title)
		}
	})
}

func TestServiceGetTasksInvalidStatusFilter(t *testing.T) {
	service := newTestService()

	_, err := service.GetTasks(TaskFilter{
		Status: "finished",
	})

	if !errors.Is(err, ErrInvalidStatus) {
		t.Errorf("expected error %v, got %v", ErrInvalidStatus, err)
	}
}
