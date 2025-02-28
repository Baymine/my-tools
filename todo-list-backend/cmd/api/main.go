package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joy_project/todo-list-backend/internal/auth"
	"github.com/joy_project/todo-list-backend/internal/database"
	"github.com/joy_project/todo-list-backend/internal/middleware"
	"github.com/joy_project/todo-list-backend/internal/models"
	"github.com/joy_project/todo-list-backend/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "TODO API: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	database.InitDB()

	// 公共路由
	http.HandleFunc("/register", middleware.CORS(logRequest(handleRegister)))
	http.HandleFunc("/login", middleware.CORS(logRequest(handleLogin)))

	// 需要认证的路由
	http.HandleFunc("/todos", middleware.CORS(logRequest(authHandler(handleTodos))))
	http.HandleFunc("/todos/", middleware.CORS(logRequest(authHandler(handleTodo))))

	logger.Println("Server starting on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func logRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}

// 认证处理器包装器
func authHandler(next http.HandlerFunc) http.HandlerFunc {
	return middleware.Auth(next)
}

// 用户注册处理
func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Printf("Error decoding register request: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	errors := validator.ValidateRegister(req)
	if len(errors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors)
		return
	}

	userID, err := database.CreateUser(req)
	if err != nil {
		logger.Printf("Error creating user: %v", err)
		if strings.Contains(err.Error(), "already exists") {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	user, err := database.GetUserByID(int(userID))
	if err != nil {
		logger.Printf("Error getting user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	token, err := auth.GenerateToken(user)
	if err != nil {
		logger.Printf("Error generating token: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Token:     token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// 用户登录处理
func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Printf("Error decoding login request: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	errors := validator.ValidateLogin(req)
	if len(errors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors)
		return
	}

	user, err := database.GetUserByEmail(req.Email)
	if err != nil {
		logger.Printf("Error getting user: %v", err)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		logger.Printf("Invalid password: %v", err)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(user)
	if err != nil {
		logger.Printf("Error generating token: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Token:     token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleTodos(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// 检查是否请求分页
		if r.URL.Query().Get("page") != "" {
			getTodosWithPagination(w, r, userID)
		} else {
			getTodos(w, r, userID)
		}
	case http.MethodPost:
		createTodo(w, r, userID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleTodo(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodPut:
		updateTodo(w, r, userID)
	case http.MethodDelete:
		deleteTodo(w, r, userID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getTodos(w http.ResponseWriter, r *http.Request, userID int) {
	todos, err := database.GetAllTodos(userID)
	if err != nil {
		logger.Printf("Error getting todos: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Printf("Retrieved %d todos for user %d", len(todos), userID)
	for _, todo := range todos {
		logger.Printf("Todo: ID=%d, Title=%s, Completed=%v, Priority=%s", todo.ID, todo.Title, todo.Completed, todo.Priority)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func getTodosWithPagination(w http.ResponseWriter, r *http.Request, userID int) {
	// 解析分页参数
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	page := 1
	pageSize := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	// 获取分页数据
	todos, total, err := database.GetTodosWithPagination(userID, page, pageSize)
	if err != nil {
		logger.Printf("Error getting todos with pagination: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Printf("Retrieved %d todos (page %d of %d) for user %d", len(todos), page, (total+pageSize-1)/pageSize, userID)

	// 构建响应
	response := map[string]interface{}{
		"todos": todos,
		"pagination": map[string]int{
			"page":       page,
			"pageSize":   pageSize,
			"total":      total,
			"totalPages": (total + pageSize - 1) / pageSize,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func createTodo(w http.ResponseWriter, r *http.Request, userID int) {
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

	todo.UserID = userID
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

func updateTodo(w http.ResponseWriter, r *http.Request, userID int) {
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
	todo.UserID = userID
	err = database.UpdateTodo(todo)
	if err != nil {
		logger.Printf("Error updating todo: %v", err)
		if strings.Contains(err.Error(), "not found or not owned") {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deleteTodo(w http.ResponseWriter, r *http.Request, userID int) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/todos/"))
	if err != nil {
		logger.Printf("Invalid todo ID: %v", err)
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	err = database.DeleteTodo(id, userID)
	if err != nil {
		logger.Printf("Error deleting todo: %v", err)
		if strings.Contains(err.Error(), "not found or not owned") {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
