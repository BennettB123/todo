package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

const ConfigDirName string = ".todo"
const DBFileName string = "data.sqlite3"

// Thin wrapper around a *sql.DB
type Database struct {
	db *sql.DB
}

func GetOrCreateDatabase() Database {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("unable to find user's home directory")
		fmt.Printf("\terror: %s", err)
		os.Exit(1)
	}

	dbDirPath := filepath.Join(userHomeDir, ConfigDirName)
	CreateDirectory(dbDirPath)

	db, err := sql.Open("sqlite", filepath.Join(dbDirPath, DBFileName))
	if err != nil {
		fmt.Println("unable to open sqlite3 database file.")
		fmt.Printf("\terror: %s", err)
		os.Exit(1)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("unable to open sqlite3 database file.")
		fmt.Printf("\terror: %s", err)
		os.Exit(1)
	}

	return Database{db}
}

func (todoDB *Database) Init() {
	_, err := todoDB.db.ExecContext(context.Background(),
		`CREATE TABLE IF NOT EXISTS todo (
		id INTEGER PRIMARY KEY,
		name string,
		status string)`)
	if err != nil {
		fmt.Println("unable to create 'todo' table in sqlite3 database file.")
		fmt.Printf("\terror: %s", err)
		os.Exit(1)
	}
}

func (database *Database) CreateTodoEntry(todo Todo) {
	_, err := database.db.ExecContext(context.Background(),
		`INSERT INTO todo (name, status) VALUES (?, ?)`,
		todo.name, todo.status)
	if err != nil {
		fmt.Println("unable to insert todo entry into database.")
		fmt.Printf("\terror: %s", err)
		os.Exit(1)
	}
}

func (database *Database) GetAllTodoEntries() []Todo {
	rows, err := database.db.QueryContext(context.Background(),
		`SELECT id, name, status FROM todo`)
	if err != nil {
		fmt.Println("unable to query todo entries from database.")
		fmt.Printf("\terror: %s", err)
		os.Exit(1)
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.id, &todo.name, &todo.status)
		if err != nil {
			fmt.Println("unable to scan todo entry from database.")
			fmt.Printf("\terror: %s", err)
			os.Exit(1)
		}
		todos = append(todos, todo)
	}
	if rows.Err() != nil {
		fmt.Println("error while iterating over todo entries from database.")
		fmt.Printf("\terror: %s", err)
		os.Exit(1)
	}

	return todos
}

func (database *Database) MarkAsDone(id uint32) {
	_, err := database.db.ExecContext(context.Background(),
		`UPDATE todo SET status = ? WHERE id = ?`, Done, id)
	if err != nil {
		fmt.Println("unable to update todo entry in database.")
		fmt.Printf("\terror: %s", err)
		os.Exit(1)
	}
}

func (database *Database) Close() {
	err := database.db.Close()
	if err != nil {
		fmt.Println("unable to close database.")
		fmt.Printf("\terror: %s", err)
	}
}
