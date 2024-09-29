package repositories

import (
	"github.com/s5k/go-tasktracker/entities"
)

type TaskRepository interface {
	Create(task *entities.Task) error
	Update(task *entities.Task) error
	Delete(id uint) error
	Get(id uint) (*entities.Task, error)
	GetAll(status string) ([]*entities.Task, error)
}
