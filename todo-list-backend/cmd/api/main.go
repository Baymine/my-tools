package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/yourusername/todo-list-backend/internal/database"
	"github.com/yourusername/todo-list-backend/internal/middleware"
	"github.com/yourusername/todo-list-backend/internal/models"
	"github.com/yourusername/todo-list-backend/internal/validator"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "TODO API: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	database.InitDB()

	http.HandleFunc("/todos", middleware.CORS(logRequest(handleTodos)))
	http.HandleFunc("/todos/", middleware.CORS(logRequest(handleTodo)))

	logger.Println("Server starting on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func logRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}

func handleTodos(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTodos(w, r)
	case http.MethodPost:
		createTodo(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleTodo(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		updateTodo(w, r)
	case http.MethodDelete:
		deleteTodo(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := database.GetAllTodos()
	if err != nil {
		logger.Printf("Error getting todos: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		logger.Printf("Error decoding todo: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	errors := validator.ValidateTodo(todo)
	if len(errors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors)
		return
	}

	id, err := database.CreateTodo(todo)
	if err != nil {
		logger.Printf("Error creating todo: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/todos/"))
	if err != nil {
		logger.Printf("Invalid todo ID: %v", err)
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	var todo models.Todo
	err = json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		logger.Printf("Error decoding todo: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	errors := validator.ValidateTodo(todo)
	if len(errors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors)
		return
	}

	todo.ID = id
	err = database.UpdateTodo(todo)
	if err != nil {
		logger.Printf("Error updating todo: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/todos/"))
	if err != nil {
		logger.Printf("Invalid todo ID: %v", err)
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	err = database.DeleteTodo(id)
	if err != nil {
		logger.Printf("Error deleting todo: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
