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

func (s *Service) GetTaskByID(id int) (Task, error) {

	task, found := s.repo.GetByID(id)
	if !found {
		return Task{}, errors.New("task not found")
	}
	return task, nil
}

func (s *Service) DeleteTask(id int) error {
	deleted := s.repo.Remove(id)
	if !deleted {
		return errors.New("task not found")
	}
	return nil
}

func (s *Service) UpdateTask(id int, req UpdateTaskRequest) (Task, error) {
	oldTask, err := s.GetTaskByID(id)
	if err != nil {
		return Task{}, err
	}

	if strings.TrimSpace(req.Title) == "" {
		return Task{}, errors.New("title is required")
	}

	if strings.TrimSpace(req.Status) == "" {
		req.Status = "Todo"
	}

	if req.Status != "todo" && req.Status != "in_progress" && req.Status != "done" {
		return Task{}, errors.New("invalid status")
	}

	oldTask.Title = strings.TrimSpace(req.Title)
	oldTask.Description = strings.TrimSpace(req.Description)
	oldTask.Status = strings.TrimSpace(req.Status)
	oldTask.Category = strings.TrimSpace(req.Category)
	oldTask.UpdatedAt = time.Now().Format(time.RFC3339)

	updatedTask, found := s.repo.UpdateTask(id, oldTask)
	if !found {
		return Task{}, errors.New("task not found")
	}

	return updatedTask, nil

}
