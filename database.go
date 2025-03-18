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

func (database *Database) Close() {
	err := database.db.Close()
	if err != nil {
		fmt.Println("unable to close database.")
		fmt.Printf("\terror: %s", err)
	}
}
