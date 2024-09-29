package services

import (
	"github.com/s5k/go-tasktracker/entities"
	"github.com/s5k/go-tasktracker/interfaces"
	repositories "github.com/s5k/go-tasktracker/repositories"
)

type TaskService struct {
	taskRepository repositories.TaskRepository
}

func NewTaskService(taskRepository repositories.TaskRepository) interfaces.TaskService {
	return &TaskService{taskRepository: taskRepository}
}

func (s *TaskService) Create(task *entities.Task) error {
	return s.taskRepository.Create(task)
}

func (s *TaskService) Update(task *entities.Task) error {
	return s.taskRepository.Update(task)
}

func (s *TaskService) Delete(id uint) error {
	return s.taskRepository.Delete(id)
}

func (s *TaskService) Get(id uint) (*entities.Task, error) {
	return s.taskRepository.Get(id)
}

func (s *TaskService) GetAll(status string) ([]*entities.Task, error) {
	return s.taskRepository.GetAll(status)
}
