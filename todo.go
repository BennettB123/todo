package main

// Todo Statuses
const (
	Open   = "OPEN"
	Closed = "CLOSED"
)

type Todo struct {
	id     uint32
	name   string
	status string
}
