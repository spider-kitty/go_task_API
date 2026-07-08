package task

import (
	"errors"
	"strings"
	"time"
)

type Service struct {
	repo *MemoryRepository
}

func NewService(repo *MemoryRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateTask(req CreateTaskRequest) (Task, error) {
	title := strings.TrimSpace(req.Title)

	if title == "" {
		return Task{}, errors.New("title is required")
	}

	status := strings.TrimSpace(req.Status)

	if status == "" {
		status = "todo"
	}

	if status != "todo" && status != "in_progress" && status != "done" {
		return Task{}, errors.New("invalid status")
	}

	now := time.Now().Format(time.RFC3339)

	task := Task{
		Title:       title,
		Description: strings.TrimSpace(req.Description),
		Status:      status,
		Category:    strings.TrimSpace(req.Category),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	createdTask := s.repo.Create(task)

	return createdTask, nil
}

func (s *Service) GetTasks() []Task {
	return s.repo.GetAll()
}
