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

func GetOrCreateDatabase(logger Logger) (Database, error) {
	logger.LogDebug("getting user's home directory.")
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return Database{}, fmt.Errorf("unable to get user's home directory: %w", err)
	}
	logger.LogDebug(fmt.Sprintf("user's home directory: %s", userHomeDir))

	dbDirPath := filepath.Join(userHomeDir, ConfigDirName)
	if err = CreateDirectory(dbDirPath, logger); err != nil {
		return Database{}, fmt.Errorf("unable to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite", filepath.Join(dbDirPath, DBFileName))
	if err != nil {
		return Database{}, fmt.Errorf("unable to open sqlite3 database file: %w", err)
	}

	if err = db.Ping(); err != nil {
		return Database{}, fmt.Errorf("unable to contact sqlite3 database: %w", err)
	}

	return Database{db, logger}, nil
}

func (database *Database) Init() error {
	_, err := database.db.ExecContext(context.Background(),
		`CREATE TABLE IF NOT EXISTS todo (
		id INTEGER PRIMARY KEY,
		name string,
		status string)`)
	if err != nil {
		return fmt.Errorf("unable to create 'todo' table in sqlite3 database file: %w", err)
	}

	return nil
}

func (database *Database) CreateTodoEntry(todo Todo) error {
	_, err := database.db.ExecContext(context.Background(),
		`INSERT INTO todo (name, status) VALUES (?, ?)`,
		todo.name, todo.status)
	if err != nil {
		return fmt.Errorf("unable to insert todo entry into database: %w", err)
	}

	return nil
}

func (database *Database) GetAllTodoEntries() ([]Todo, error) {
	rows, err := database.db.QueryContext(context.Background(),
		`SELECT id, name, status FROM todo`)
	if err != nil {
		return nil, fmt.Errorf("unable to query todo entries from database: %w", err)
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.id, &todo.name, &todo.status); err != nil {
			return nil, fmt.Errorf("unable to scan todo entry from database: %w", err)
		}
		todos = append(todos, todo)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("error while iterating over todo entries from database: %w", err)
	}

	return todos, nil
}

func (database *Database) ChangeTodoStatus(id uint32, status string) error {
	_, err := database.db.ExecContext(context.Background(),
		`UPDATE todo SET status = ? WHERE id = ?`, status, id)
	if err != nil {
		return fmt.Errorf("unable to update todo status in database: %w", err)
	}

	return nil
}

func (database *Database) ChangeName(id uint32, name string) error {
	_, err := database.db.ExecContext(context.Background(),
		`UPDATE todo SET name = ? WHERE id = ?`, name, id)
	if err != nil {
		return fmt.Errorf("unable to update todo name in database: %w", err)
	}

	return nil
}

func (database *Database) DeleteEntry(id uint32) error {
	_, err := database.db.ExecContext(context.Background(),
		`DELETE FROM todo WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("unable to delete todo entry from database: %w", err)
	}

	return nil
}

func (database *Database) Close() error {
	if err := database.db.Close(); err != nil {
		return fmt.Errorf("unable to close database: %w", err)
	}

	return nil
}
