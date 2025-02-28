package database

import (
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joy_project/todo-list-backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/todo_list?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to the database successfully")
}

// 用户相关操作

// CreateUser 创建新用户
func CreateUser(user models.RegisterRequest) (int64, error) {
	// 检查邮箱是否已存在
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", user.Email).Scan(&count)
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, errors.New("email already exists")
	}

	// 检查用户名是否已存在
	err = DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", user.Username).Scan(&count)
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, errors.New("username already exists")
	}

	// 哈希密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	// 插入用户
	result, err := DB.Exec(
		"INSERT INTO users (username, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		user.Username, user.Email, hashedPassword, time.Now(), time.Now(),
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// GetUserByEmail 通过邮箱获取用户
func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := DB.QueryRow(
		"SELECT id, username, email, password, created_at, updated_at FROM users WHERE email = ?",
		email,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// GetUserByID 通过ID获取用户
func GetUserByID(id int) (models.User, error) {
	var user models.User
	err := DB.QueryRow(
		"SELECT id, username, email, password, created_at, updated_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// Todo相关操作

// GetAllTodos 获取指定用户的所有待办事项
func GetAllTodos(userID int) ([]models.Todo, error) {
	rows, err := DB.Query("SELECT id, title, completed, priority, user_id, created_at, updated_at FROM todos WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.Priority, &todo.UserID, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	if len(todos) == 0 {
		return []models.Todo{}, nil
	}

	return todos, nil
}

// CreateTodo 创建待办事项
func CreateTodo(todo models.Todo) (int64, error) {
	result, err := DB.Exec("INSERT INTO todos (title, completed, priority, user_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		todo.Title, todo.Completed, todo.Priority, todo.UserID, time.Now(), time.Now())
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// UpdateTodo 更新待办事项
func UpdateTodo(todo models.Todo) error {
	// 验证待办事项属于当前用户
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM todos WHERE id = ? AND user_id = ?", todo.ID, todo.UserID).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("todo not found or not owned by user")
	}

	_, err = DB.Exec("UPDATE todos SET title = ?, completed = ?, priority = ?, updated_at = ? WHERE id = ?",
		todo.Title, todo.Completed, todo.Priority, time.Now(), todo.ID)
	return err
}

// DeleteTodo 删除待办事项
func DeleteTodo(id int, userID int) error {
	// 验证待办事项属于当前用户
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM todos WHERE id = ? AND user_id = ?", id, userID).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("todo not found or not owned by user")
	}

	_, err = DB.Exec("DELETE FROM todos WHERE id = ?", id)
	return err
}
