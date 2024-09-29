package file

// connect to a file and store the tasks in it
// use json to store the tasks

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/s5k/go-tasktracker/entities"
	repositories "github.com/s5k/go-tasktracker/repositories"
)

type TaskRepository struct {
	file *os.File
}

func NewTaskRepository(file *os.File) repositories.TaskRepository {
	return &TaskRepository{file: file}
}

func (r *TaskRepository) restructureFile() error {
	return nil
}

func (r *TaskRepository) Create(task *entities.Task) error {
	// Read existing tasks
	var tasks []*entities.Task
	fileContent, err := os.ReadFile(r.file.Name())
	if err != nil {
		return err
	}

	if len(fileContent) > 0 {
		if err := json.Unmarshal(fileContent, &tasks); err != nil {
			return err
		}
	}

	// Generate a new ID
	newID := 1
	if len(tasks) > 0 {
		lastTask := tasks[len(tasks)-1]
		lastID := lastTask.ID
		if err != nil {
			return err
		}
		newID = int(lastID) + 1
	}
	task.ID = uint(newID)

	// Add the new task
	tasks = append(tasks, task)

	// Write all tasks back to the file
	updatedContent, err := json.Marshal(tasks)
	if err != nil {
		return err
	}

	if err := os.WriteFile(r.file.Name(), updatedContent, 0644); err != nil {
		return err
	}

	return nil
}
func (r *TaskRepository) Update(task *entities.Task) error {
	if task.ID == 0 {
		return errors.New("task ID is required")
	}

	// Read existing tasks
	var tasks []*entities.Task
	fileContent, err := os.ReadFile(r.file.Name())
	if err != nil {
		return err
	}

	if len(fileContent) > 0 {
		if err := json.Unmarshal(fileContent, &tasks); err != nil {
			return err
		}
	}

	// Update the task
	for i, t := range tasks {
		if t.ID == task.ID {
			tasks[i] = task
			break
		}
	}

	// Write all tasks back to the file
	updatedContent, err := json.Marshal(tasks)
	if err != nil {
		return err
	}

	if err := os.WriteFile(r.file.Name(), updatedContent, 0644); err != nil {
		return err
	}

	return nil
}

func (r *TaskRepository) Delete(id uint) error {
	if id == 0 {
		return errors.New("task ID is required")
	}

	if err := r.restructureFile(); err != nil {
		return err
	}

	// Read existing tasks
	var tasks []*entities.Task
	fileContent, err := os.ReadFile(r.file.Name())
	if err != nil {
		return err
	}

	if len(fileContent) > 0 {
		if err := json.Unmarshal(fileContent, &tasks); err != nil {
			return err
		}
	}

	// Delete the task
	for i, t := range tasks {
		if t.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			break
		}
	}

	// Write all tasks back to the file
	updatedContent, err := json.Marshal(tasks)
	if err != nil {
		return err
	}

	if err := os.WriteFile(r.file.Name(), updatedContent, 0644); err != nil {
		return err
	}

	return nil
}

func (r *TaskRepository) Get(id uint) (*entities.Task, error) {
	if id == 0 {
		return nil, errors.New("task ID is required")
	}

	// Read existing tasks
	var tasks []*entities.Task
	fileContent, err := os.ReadFile(r.file.Name())
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return nil, err
	}

	if len(fileContent) > 0 {
		if err := json.Unmarshal(fileContent, &tasks); err != nil {
			fmt.Printf("Error unmarshalling file: %v\n", err)
			return nil, err
		}
	}

	for _, t := range tasks {
		if t.ID == id {
			return t, nil
		}
	}

	return nil, errors.New("task not found")
}

func (r *TaskRepository) GetAll(status string) ([]*entities.Task, error) {
	// Read existing tasks
	var tasks []*entities.Task
	fileContent, err := os.ReadFile(r.file.Name())
	if err != nil {
		return nil, err
	}

	if len(fileContent) > 0 {
		if err := json.Unmarshal(fileContent, &tasks); err != nil {
			return nil, err
		}
	}

	if status != "" {
		//get all tasks with the status
		filteredTasks := make([]*entities.Task, 0)
		for _, t := range tasks {
			if t.Status == status {
				filteredTasks = append(filteredTasks, t)
			}
		}
		return filteredTasks, nil
	}

	return tasks, nil
}

func (r *TaskRepository) Close() error {
	return r.file.Close()
}
