package main

import (
	"fmt"
	"strings"

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
	return PrintTodos(context)
}

type NewCmd struct {
	Name string `arg:"" help:"The name of the TODO entry."`
}

func (n *NewCmd) Run(context Context) error {
	context.logger.LogDebug(fmt.Sprintf("sanitizing provided TODO name '%s'", n.Name))
	name := SanitizeName(n.Name)
	todo := NewTodo(name)

	context.logger.LogDebug(fmt.Sprintf("creating new TODO entry with name [%s] and status [%s]", todo.name, todo.status))
	err := context.db.CreateTodoEntry(todo)
	if err != nil {
		return err
	}

	err = PrintTodos(context)
	return err
}

type DoneCmd struct {
	Ids []uint32 `arg:"" help:"IDs of TODO entries to mark Done. Multiple, space-separated values are supported."`
}

func (d *DoneCmd) Run(context Context) error {
	context.logger.LogDebug(fmt.Sprintf("Marking TODO entries as Done: %v", d.Ids))
	for _, id := range d.Ids {
		err := context.db.ChangeTodoStatus(id, Done)
		if err != nil {
			context.logger.LogError(fmt.Sprintf("unable to mark entry with ID '%d' as Done: %v", id, err))
		}
	}

	err := PrintTodos(context)
	return err
}

type OpenCmd struct {
	Ids []uint32 `arg:"" help:"IDs of TODO entries to mark Open. Multiple, space-separated values are supported."`
}

func (d *OpenCmd) Run(context Context) error {
	context.logger.LogDebug(fmt.Sprintf("Marking TODO entries as Open: %v", d.Ids))
	for _, id := range d.Ids {
		err := context.db.ChangeTodoStatus(id, Open)
		if err != nil {
			context.logger.LogError(fmt.Sprintf("unable to mark entry with ID '%d' as Open: %v", id, err))
		}
	}

	err := PrintTodos(context)
	return err
}

type EditCmd struct {
	Id   uint32 `arg:"" help:"ID of TODO entry to edit."`
	Name string `arg:"" help:"New name for the TODO entry."`
}

func (e *EditCmd) Run(context Context) error {
	context.logger.LogDebug(fmt.Sprintf("sanitizing provided TODO name '%s'", e.Name))
	name := SanitizeName(e.Name)

	context.logger.LogDebug(fmt.Sprintf("editing name of entry with ID '%d' to [%s]", e.Id, name))
	err := context.db.ChangeName(e.Id, name)
	if err != nil {
		context.logger.LogError(fmt.Sprintf("unable to change name of entry with ID '%d': %v", e.Id, err))
	}

	err = PrintTodos(context)
	return err
}

var CLI struct {
	List  ListCmd `cmd:"" default:"1" aliases:"ls" help:"List TODO entries."`
	New   NewCmd  `cmd:"" aliases:"n" help:"Create a new TODO entry."`
	Done  DoneCmd `cmd:"" aliases:"d" help:"Mark existing TODO entries as Done."`
	Open  OpenCmd `cmd:"" aliases:"o" help:"Mark existing TODO entries as Open."`
	Edit  EditCmd `cmd:"" aliases:"e" help:"Edit the name of an existing TODO entry."`
	Debug bool    `help:"Enable debug mode for verbose logging."`
}

func PrintTodos(context Context) error {
	context.logger.LogDebug("retrieving all TODO entries")
	todos, err := context.db.GetAllTodoEntries()
	if err != nil {
		return err
	}

	context.logger.LogDebug(fmt.Sprintf("%d TODO entries found", len(todos)))

	for _, todo := range todos {
		if todo.status == Open {
			fmt.Printf("%d: [ ] %s\n", todo.id, todo.name)
		}
	}

	for _, todo := range todos {
		if todo.status == Done {
			fmt.Printf("%d: [X] %s\n", todo.id, todo.name)
		}
	}

	return nil
}

func main() {
	ctx := kong.Parse(&CLI,
		kong.Name("todo"),
		kong.Description("A simple command-line TODO list manager."))

	logger := Logger{debug: CLI.Debug}
	if CLI.Debug {
		logger.LogDebug("Debug mode enabled.")
	}

	db, err := GetOrCreateDatabase(logger)
	if err != nil {
		logger.LogError(err.Error())
		return
	}

	defer db.Close()
	err = db.Init()
	if err != nil {

		logger.LogError(err.Error())
		return
	}

	err = ctx.Run(Context{db, logger})
	if err != nil {
		logger.LogError(err.Error())
	}
}

func SanitizeName(name string) string {
	return strings.Split(strings.ReplaceAll(name, "\r\n", "\n"), "\n")[0]
}
