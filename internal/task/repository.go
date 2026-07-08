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

func (r *MemoryRepository) GetByID(id int) (Task, bool) {

	if id == 0 || id > len(r.tasks) {
		return Task{}, false
	}

	for _, task := range r.tasks {
		if task.ID == id {
			return task, true
		}
	}

	return Task{}, false
}

func (r *MemoryRepository) Remove(id int) bool {
	for i, task := range r.tasks {
		if task.ID == id {
			r.tasks = append(r.tasks[:i], r.tasks[i+1:]...)
			return true
		}
	}
	return false
}

func (r *MemoryRepository) UpdateTask(id int, task Task) (Task, bool) {

	for i, t := range r.tasks {
		if t.ID == id {
			r.tasks[i] = task
			return task, true
		}
	}
	return Task{}, false
}
