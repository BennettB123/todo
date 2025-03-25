package main

import "fmt"

// Todo Statuses
const (
	Open = "OPEN"
	Done = "DONE"
)

type Todo struct {
	id       uint32
	name     string
	status   string
	archived bool
}

// NewTodo creates a new Todo instance.
//   The Id is set to 0 as a placeholder, since the database will assign the Id.
//   The status is set to Open.
func NewTodo(name string) Todo {
	return Todo{0, name, Open, false}
}

func (t Todo) String() string {
	checkbox := "[ ]"
	if t.status == Done {
		checkbox = "[X]"
	}

	return fmt.Sprintf("%-3d %s %s", t.id, checkbox, t.name)
}
