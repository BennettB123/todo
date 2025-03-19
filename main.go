package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	_ "modernc.org/sqlite"
)

type Context struct {
	db      Database
	verbose bool
}

type ListCmd struct {
}

func (l *ListCmd) Run(context Context) error {
	todos := context.db.GetAllTodoEntries()
	for _, todo := range todos {
		fmt.Printf("%d: %s - %s\n", todo.id, todo.status, todo.name)
	}

	return nil
}

type NewCmd struct {
	// TODO: make Name a required/positional field (always first arg after `new`)
	Name string `short:"n" help:"The name of the TODO entry."`
}

func (n *NewCmd) Run(context Context) error {
	if n.Name == "" {
		return fmt.Errorf("name is required")
	}

	todo := NewTodo(n.Name)
	context.db.CreateTodoEntry(todo)

	return nil
}

type DoneCmd struct {
	Id uint32 `short:"i" help:"The id of the TODO entry to mark as done."`
}

func (d *DoneCmd) Run(context Context) error {
	context.db.MarkAsDone(d.Id)

	return nil
}

var CLI struct {
	List ListCmd `cmd:"" help:"List TODO entries."`
	New  NewCmd  `cmd:"" help:"Create a new TODO entry."`
	Done DoneCmd `cmd:"" help:"Mark an existing TODO entry as Done."`
}

func main() {
	db := GetOrCreateDatabase()
	defer db.Close()
	db.Init()

	ctx := kong.Parse(&CLI)
	ctx.Run(Context{db: db, verbose: false})
}
