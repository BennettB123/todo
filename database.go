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
	db     *sql.DB
	logger Logger
}

func GetOrCreateDatabase(logger Logger) Database {
	logger.Log(Debug, "getting user's home directory.")
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		logger.LogError("unable to find user's home directory.", err)
		os.Exit(1)
	}
	logger.Log(Debug, fmt.Sprintf("user's home directory: %s", userHomeDir))

	dbDirPath := filepath.Join(userHomeDir, ConfigDirName)
	if err = CreateDirectory(dbDirPath, logger); err != nil {
		logger.LogError("unable to create database directory.", err)
	}

	db, err := sql.Open("sqlite", filepath.Join(dbDirPath, DBFileName))
	if err != nil {
		logger.LogError("unable to open sqlite3 database file.", err)
		os.Exit(1)
	}

	if err = db.Ping(); err != nil {
		logger.LogError("unable to open sqlite3 database file.", err)
		os.Exit(1)
	}

	return Database{db, logger}
}

func (database *Database) Init() {
	_, err := database.db.ExecContext(context.Background(),
		`CREATE TABLE IF NOT EXISTS todo (
		id INTEGER PRIMARY KEY,
		name string,
		status string)`)
	if err != nil {
		database.logger.LogError("unable to create 'todo' table in sqlite3 database file.", err)
		os.Exit(1)
	}
}

func (database *Database) CreateTodoEntry(todo Todo) {
	_, err := database.db.ExecContext(context.Background(),
		`INSERT INTO todo (name, status) VALUES (?, ?)`,
		todo.name, todo.status)
	if err != nil {
		database.logger.LogError("unable to insert todo entry into database.", err)
		os.Exit(1)
	}
}

func (database *Database) GetAllTodoEntries() []Todo {
	rows, err := database.db.QueryContext(context.Background(),
		`SELECT id, name, status FROM todo`)
	if err != nil {
		database.logger.LogError("unable to query todo entries from database.", err)
		os.Exit(1)
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.id, &todo.name, &todo.status); err != nil {
			database.logger.LogError("unable to scan todo entry from database.", err)
			os.Exit(1)
		}
		todos = append(todos, todo)
	}
	if rows.Err() != nil {
		database.logger.LogError("error while iterating over todo entries from database.", err)
		os.Exit(1)
	}

	return todos
}

func (database *Database) MarkAsDone(id uint32) {
	_, err := database.db.ExecContext(context.Background(),
		`UPDATE todo SET status = ? WHERE id = ?`, Done, id)
	if err != nil {
		database.logger.LogError("unable to update todo entry in database.", err)
		os.Exit(1)
	}
}

func (database *Database) Close() {
	if err := database.db.Close(); err != nil {
		database.logger.LogError("unable to close database.", err)
	}
}
