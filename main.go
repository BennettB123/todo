package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	_ "modernc.org/sqlite"
)

type Context struct {
	db     Database
	logger Logger
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
	Name string `arg:"" help:"The name of the TODO entry."`
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
	Ids []uint32 `arg:"" help:"IDs of TODO entries to mark Done. Multiple, space-separated values are supported."`
}

func (d *DoneCmd) Run(context Context) error {
	context.logger.Log(Debug, fmt.Sprintf("Marking TODO entries as Done: %v", d.Ids))
	for _, id := range d.Ids {
		context.db.ChangeTodoStatus(id, Done)
	}

	return nil
}

type OpenCmd struct {
	Ids []uint32 `arg:"" help:"IDs of TODO entries to mark Open. Multiple, space-separated values are supported."`
}

func (d *OpenCmd) Run(context Context) error {
	context.logger.Log(Debug, fmt.Sprintf("Marking TODO entries as Open: %v", d.Ids))
	for _, id := range d.Ids {
		context.db.ChangeTodoStatus(id, Open)
	}

	return nil
}

var CLI struct {
	List  ListCmd `cmd:"" help:"List TODO entries."`
	New   NewCmd  `cmd:"" help:"Create a new TODO entry."`
	Done  DoneCmd `cmd:"" help:"Mark existing TODO entries as Done."`
	Open  OpenCmd `cmd:"" help:"Mark existing TODO entries as Open."`
	Debug bool    `help:"Enable debug mode for verbose logging."`
}

func main() {
	ctx := kong.Parse(&CLI)

	logger := Logger{debug: CLI.Debug}
	if CLI.Debug {
		logger.Log(Debug, "Debug mode enabled.")
	}

	db := GetOrCreateDatabase(logger)
	defer db.Close()
	db.Init()

	ctx.Run(Context{db, logger})
}
