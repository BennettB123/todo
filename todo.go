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

// NewTodo creates a new Todo instance.
//   The Id is set to 0 as a placeholder, since the database will assign the Id.
func NewTodo(name string) Todo {
	return Todo{0, name, Open}
}
