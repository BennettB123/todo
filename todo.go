package main

// Todo Statuses
const (
	Open   = "OPEN"
	Closed = "CLOSED"
)

type Todo struct {
	name   string
	status string
}