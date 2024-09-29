package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/s5k/go-tasktracker/entities"
	"github.com/s5k/go-tasktracker/infrastructure/db/file"
	"github.com/s5k/go-tasktracker/services"
)

func printHelp() {
	fmt.Println("Usage: tasktracker <command> <args>")
	fmt.Println("Available commands:")
	fmt.Println("  - add <task_description>")
	fmt.Println("  - list")
	fmt.Println("  - list <status>")
	fmt.Println("  - update <task_id> <new_description>")
	fmt.Println("  - delete <task_id>")
	fmt.Println("  - mark-in-progress <task_id>")
	fmt.Println("  - mark-done <task_id>")
	fmt.Println("  - help")
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("No command provided.")

		printHelp()
		return
	}
	command := args[0]

	taskFile, err := os.OpenFile("./data/tasks.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Error opening task file: %v", err)
	}
	defer taskFile.Close()

	taskRepository := file.NewTaskRepository(taskFile)
	taskService := services.NewTaskService(taskRepository)

	switch command {
	case "add":
		if len(args) < 2 {
			fmt.Println("No task description provided.")
			return
		}

		taskDescription := args[1]

		task := entities.Task{
			Description: taskDescription,
			Status:      "todo",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		err = taskService.Create(&task)
		if err != nil {
			log.Fatalf("Error creating task: %v", err)
		}
		log.Printf("Task created: %v", task)
	case "list":
		status := ""
		if len(args) > 1 {
			status = args[1]
		}

		tasks, err := taskService.GetAll(status)
		if err != nil {
			log.Fatalf("Error getting tasks: %v", err)
		}
		for _, task := range tasks {
			fmt.Printf("[id: %d, description: %s, status: %s, created_at: %s, updated_at: %s]\n", task.ID, task.Description, task.Status, task.CreatedAt, task.UpdatedAt)
		}
	case "update":
		if len(args) < 3 {
			fmt.Println("No task id or new description provided.")
			return
		}

		taskID, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatalf("Error converting task id to int: %v", err)
		}
		newDescription := args[2]

		task, err := taskService.Get(uint(taskID))
		if err != nil {
			log.Fatalf("Error getting task: %v", err)
		}
		task.Description = newDescription
		err = taskService.Update(task)
		if err != nil {
			log.Fatalf("Error updating task: %v", err)
		}
		log.Printf("Task updated: %v", task)
	case "delete":
		if len(args) < 2 {
			fmt.Println("No task id provided.")
			return
		}

		taskID, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatalf("Error converting task id to int: %v", err)
		}
		fmt.Println("Deleting task:", args[1])
		err = taskService.Delete(uint(taskID))
		if err != nil {
			log.Fatalf("Error deleting task: %v", err)
		}
		log.Printf("Task deleted: %v", taskID)
	case "mark-in-progress":
		if len(args) < 2 {
			fmt.Println("No task id provided.")
			return
		}

		taskID, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatalf("Error converting task id to int: %v", err)
			return
		}
		task, err := taskService.Get(uint(taskID))
		if err != nil {
			log.Fatalf("Error getting task: %v", err)
			return
		}
		task.Status = "in-progress"
		err = taskService.Update(task)
		if err != nil {
			log.Fatalf("Error updating task: %v", err)
		}
		log.Printf("Task updated: %v", task)
	case "mark-done":
		if len(args) < 2 {
			fmt.Println("No task id provided.")
			return
		}

		taskID, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatalf("Error converting task id to int: %v", err)
		}
		task, err := taskService.Get(uint(taskID))
		if err != nil {
			log.Fatalf("Error getting task: %v", err)
		}
		task.Status = "done"
		err = taskService.Update(task)
		if err != nil {
			log.Fatalf("Error updating task: %v", err)
		}
		log.Printf("Task updated: %v", task)
	case "help":
		printHelp()
	default:
		fmt.Println("Unknown command:", command)
		printHelp()
	}
}
