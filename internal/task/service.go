package task

import (
	"strings"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateTask(req CreateTaskRequest) (Task, error) {
	title := strings.TrimSpace(req.Title)
	description := strings.TrimSpace(req.Description)
	category := strings.TrimSpace(req.Category)

	if err := validateTitle(title); err != nil {
		return Task{}, err
	}
	if err := validateDescription(description); err != nil {
		return Task{}, err
	}
	if err := validateCategory(category); err != nil {
		return Task{}, err
	}

	status, err := normalizeStatus(req.Status)
	if err != nil {
		return Task{}, err
	}

	now := time.Now().Format(time.RFC3339)

	task := Task{
		Title:       title,
		Description: description,
		Status:      status,
		Category:    category,
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

	if filter.Status != "" {
		status, err := normalizeStatus(filter.Status)
		if err != nil {
			return []Task{}, err
		}
		filter.Status = status
	}

	if err := validateSearchFilter(filter.Search); err != nil {
		return []Task{}, err
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
		return Task{}, ErrTaskNotFound
	}
	return task, nil
}

func (s *Service) DeleteTask(id int) error {
	deleted := s.repo.Remove(id)
	if !deleted {
		return ErrTaskNotFound
	}
	return nil
}

func (s *Service) UpdateTask(id int, req UpdateTaskRequest) (Task, error) {
	oldTask, err := s.GetTaskByID(id)
	if err != nil {
		return Task{}, err
	}

	title := strings.TrimSpace(req.Title)
	description := strings.TrimSpace(req.Description)
	category := strings.TrimSpace(req.Category)

	if err := validateTitle(title); err != nil {
		return Task{}, err
	}
	if err := validateDescription(description); err != nil {
		return Task{}, err
	}
	if err := validateCategory(category); err != nil {
		return Task{}, err
	}

	status := strings.TrimSpace(req.Status)
	if status == "" {
		status = oldTask.Status
	} else {
		status, err = normalizeStatus(status)
		if err != nil {
			return Task{}, err
		}
	}

	oldTask.Title = title
	oldTask.Description = description
	oldTask.Status = status
	oldTask.Category = category
	oldTask.UpdatedAt = time.Now().Format(time.RFC3339)

	updatedTask, found := s.repo.UpdateTask(id, oldTask)
	if !found {
		return Task{}, ErrTaskNotFound
	}

	return updatedTask, nil
}
