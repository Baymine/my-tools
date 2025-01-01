package validator

import (
	"github.com/yourusername/todo-list-backend/internal/models"
)

func ValidateTodo(todo models.Todo) map[string]string {
	errors := make(map[string]string)

	if todo.Title == "" {
		errors["title"] = "Title is required"
	}

	if len(todo.Title) > 255 {
		errors["title"] = "Title must be less than 255 characters"
	}

	if todo.Priority != "low" && todo.Priority != "medium" && todo.Priority != "high" {
		errors["priority"] = "Priority must be low, medium, or high"
	}

	return errors
}
