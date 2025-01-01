package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/yourusername/todo-list-backend/internal/models"
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

func GetAllTodos() ([]models.Todo, error) {
	rows, err := DB.Query("SELECT id, title, completed, priority, created_at, updated_at FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.Priority, &todo.CreatedAt, &todo.UpdatedAt)
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

func CreateTodo(todo models.Todo) (int64, error) {
	result, err := DB.Exec("INSERT INTO todos (title, completed, priority, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		todo.Title, todo.Completed, todo.Priority, time.Now(), time.Now())
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func UpdateTodo(todo models.Todo) error {
	_, err := DB.Exec("UPDATE todos SET title = ?, completed = ?, priority = ?, updated_at = ? WHERE id = ?",
		todo.Title, todo.Completed, todo.Priority, time.Now(), todo.ID)
	return err
}

func DeleteTodo(id int) error {
	_, err := DB.Exec("DELETE FROM todos WHERE id = ?", id)
	return err
}
