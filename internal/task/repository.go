package task

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

func (r *MemoryRepository) Create(task Task) Task {
	task.ID = r.nextID
	r.nextID++

	r.tasks = append(r.tasks, task)

	return task
}

func (r *MemoryRepository) GetAll() []Task {
	return r.tasks
}
