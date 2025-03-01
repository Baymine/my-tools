package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 连接到 MySQL 服务器（添加 multiStatements=true 参数）
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/?multiStatements=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 创建数据库
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS todo_list")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("数据库 todo_list 创建成功")

	// 切换到 todo_list 数据库
	_, err = db.Exec("USE todo_list")
	if err != nil {
		log.Fatal(err)
	}

	// 删除现有的表（注意顺序：先删除有外键约束的表）
	_, err = db.Exec("DROP TABLE IF EXISTS todos")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("删除 todos 表成功")

	_, err = db.Exec("DROP TABLE IF EXISTS users")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("删除 users 表成功")

	// 创建 users 表
	createUsersTable := `
	CREATE TABLE users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(50) NOT NULL UNIQUE,
		email VARCHAR(100) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)`

	_, err = db.Exec(createUsersTable)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("users 表创建成功")

	// 创建 todos 表
	createTodosTable := `
	CREATE TABLE todos (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		completed BOOLEAN DEFAULT FALSE,
		priority ENUM('low', 'medium', 'high') DEFAULT 'medium',
		user_id INT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)`

	_, err = db.Exec(createTodosTable)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("todos 表创建成功")

	log.Println("数据库设置成功完成")
}
