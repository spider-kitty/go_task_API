package task

import (
	"encoding/json"
	"os"
)

type JSONRepository struct {
	filePath string
	tasks    []Task
	nextID   int
}

func NewJSONRepository(filePath string) (*JSONRepository, error) {
	repo := &JSONRepository{
		filePath: filePath,
		tasks:    []Task{},
		nextID:   1,
	}

	err := repo.load()
	if err != nil {
		return nil, err
	}

	repo.nextID = repo.getNextID()

	return repo, nil
}

func (r *JSONRepository) load() error {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	if len(data) == 0 {
		return nil
	}

	err = json.Unmarshal(data, &r.tasks)
	if err != nil {
		return err
	}

	return nil
}

func (r *JSONRepository) save() error {
	data, err := json.MarshalIndent(r.tasks, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(r.filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (r *JSONRepository) getNextID() int {
	maxID := 0

	for _, task := range r.tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}

	return maxID + 1
}

func (r *JSONRepository) Create(task Task) Task {
	task.ID = r.nextID
	r.nextID++

	r.tasks = append(r.tasks, task)

	err := r.save()
	if err != nil {
		return task
	}

	return task
}

func (r *JSONRepository) GetAll() []Task {
	return r.tasks
}

func (r *JSONRepository) GetByID(id int) (Task, bool) {
	for _, task := range r.tasks {
		if task.ID == id {
			return task, true
		}
	}

	return Task{}, false
}

func (r *JSONRepository) UpdateTask(id int, task Task) (Task, bool) {
	for i, t := range r.tasks {
		if t.ID == id {
			r.tasks[i] = task

			err := r.save()
			if err != nil {
				return task, true
			}

			return task, true
		}
	}

	return Task{}, false
}

func (r *JSONRepository) Remove(id int) bool {
	for i, task := range r.tasks {
		if task.ID == id {
			r.tasks = append(r.tasks[:i], r.tasks[i+1:]...)

			err := r.save()
			if err != nil {
				return true
			}

			return true
		}
	}

	return false
}
