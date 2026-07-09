package task

type Repository interface {
	Create(task Task) (Task, error)
	GetAll() ([]Task, error)
	GetByID(id int) (Task, error)
	UpdateTask(id int, task Task) (Task, error)
	Remove(id int) error
}

type MemoryRepository struct {
	tasks  []Task
	nextID int
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		tasks:  []Task{},
		nextID: 1,
	}
}
func (r *MemoryRepository) Create(task Task) (Task, error) {
	task.ID = r.nextID
	r.nextID++

	r.tasks = append(r.tasks, task)

	return task, nil
}

func (r *MemoryRepository) GetAll() ([]Task, error) {
	return r.tasks, nil
}

func (r *MemoryRepository) GetByID(id int) (Task, error) {
	for _, task := range r.tasks {
		if task.ID == id {
			return task, nil
		}
	}

	return Task{}, ErrTaskNotFound
}
func (r *MemoryRepository) Remove(id int) error {
	for i, task := range r.tasks {
		if task.ID == id {
			r.tasks = append(r.tasks[:i], r.tasks[i+1:]...)
			return nil
		}
	}

	return ErrTaskNotFound
}

func (r *MemoryRepository) UpdateTask(id int, task Task) (Task, error) {
	for i, t := range r.tasks {
		if t.ID == id {
			r.tasks[i] = task
			return task, nil
		}
	}

	return Task{}, ErrTaskNotFound
}
