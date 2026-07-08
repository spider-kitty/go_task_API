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

	if !isValidStatus(status) {
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

func (s *Service) GetTasks(filter TaskFilter) ([]Task, error) {
	filter.Status = strings.TrimSpace(filter.Status)
	filter.Category = strings.TrimSpace(filter.Category)
	filter.Search = strings.TrimSpace(filter.Search)

	if filter.Status != "" && !isValidStatus(filter.Status) {
		return nil, errors.New("invalid status filter: " + filter.Status)
	}

	tasks := s.repo.GetAll()
	var filteredTasks []Task

	search := strings.ToLower(filter.Search)
	category := strings.ToLower(filter.Category)

	for _, task := range tasks {
		if filter.Status != "" && task.Status != filter.Status {
			continue
		}
		if category != "" && strings.ToLower(task.Category) != category {
			continue
		}
		if search != "" {
			title := strings.ToLower(task.Title)
			description := strings.ToLower(task.Description)

			if !strings.Contains(title, search) && !strings.Contains(description, search) {
				continue
			}
		}
		filteredTasks = append(filteredTasks, task)
	}

	return filteredTasks, nil
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

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return Task{}, errors.New("title is required")
	}

	Status := strings.TrimSpace(req.Status)
	if Status == "" {
		req.Status = "Todo"
	}

	if !isValidStatus(Status) {
		return Task{}, errors.New("invalid status")
	}
	oldTask.Title = title
	oldTask.Description = strings.TrimSpace(req.Description)
	oldTask.Status = Status
	oldTask.Category = strings.TrimSpace(req.Category)
	oldTask.UpdatedAt = time.Now().Format(time.RFC3339)

	updatedTask, found := s.repo.UpdateTask(id, oldTask)
	if !found {
		return Task{}, errors.New("task not found")
	}

	return updatedTask, nil

}

func isValidStatus(status string) bool {
	return status == "todo" || status == "in_progress" || status == "done"
}
