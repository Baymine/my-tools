package validator

import (
	"regexp"

	"github.com/joy_project/todo-list-backend/internal/models"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

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

func ValidateRegister(req models.RegisterRequest) map[string]string {
	errors := make(map[string]string)

	// 验证用户名
	if req.Username == "" {
		errors["username"] = "Username is required"
	} else if len(req.Username) < 3 {
		errors["username"] = "Username must be at least 3 characters"
	} else if len(req.Username) > 50 {
		errors["username"] = "Username must be less than 50 characters"
	}

	// 验证邮箱
	if req.Email == "" {
		errors["email"] = "Email is required"
	} else if !emailRegex.MatchString(req.Email) {
		errors["email"] = "Invalid email format"
	}

	// 验证密码
	if req.Password == "" {
		errors["password"] = "Password is required"
	} else if len(req.Password) < 6 {
		errors["password"] = "Password must be at least 6 characters"
	}

	return errors
}

func ValidateLogin(req models.LoginRequest) map[string]string {
	errors := make(map[string]string)

	// 验证邮箱
	if req.Email == "" {
		errors["email"] = "Email is required"
	} else if !emailRegex.MatchString(req.Email) {
		errors["email"] = "Invalid email format"
	}

	// 验证密码
	if req.Password == "" {
		errors["password"] = "Password is required"
	}

	return errors
}
