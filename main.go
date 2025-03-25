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
	Archived bool `short:"a" help:"Include archived TODO entries in the list."`
}

func (l *ListCmd) Run(context Context) error {
	return PrintTodos(context, l.Archived)
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

	err = PrintTodos(context, false)
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

	err := PrintTodos(context, false)
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

	err := PrintTodos(context, false)
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

	err = PrintTodos(context, false)
	return err
}

type DeleteCmd struct {
	Ids []uint32 `arg:"" help:"IDs of TODO entries to delete. Multiple, space-separated values are supported."`
}

func (d *DeleteCmd) Run(context Context) error {
	context.logger.LogDebug(fmt.Sprintf("Deleting TODO entries: %v", d.Ids))
	for _, id := range d.Ids {
		err := context.db.DeleteEntry(id)
		if err != nil {
			context.logger.LogError(fmt.Sprintf("unable to delete entry with ID '%d': %v", id, err))
		}
	}

	err := PrintTodos(context, false)
	return err
}

type ArchiveCmd struct {
	Ids []uint32 `arg:"" help:"IDs of TODO entries to archive. Multiple, space-separated values are supported."`
}

func (a *ArchiveCmd) Run(context Context) error {
	context.logger.LogDebug(fmt.Sprintf("Archiving TODO entries: %v", a.Ids))
	for _, id := range a.Ids {
		err := context.db.ArchiveEntry(id)
		if err != nil {
			context.logger.LogError(fmt.Sprintf("unable to archive entry with ID '%d': %v", id, err))
		}
	}

	err := PrintTodos(context, false)
	return err
}

var CLI struct {
	List    ListCmd    `cmd:"" default:"1" aliases:"ls" help:"List TODO entries."`
	New     NewCmd     `cmd:"" help:"Create a new TODO entry."`
	Done    DoneCmd    `cmd:"" help:"Mark existing TODO entries as Done."`
	Open    OpenCmd    `cmd:"" help:"Mark existing TODO entries as Open."`
	Edit    EditCmd    `cmd:"" help:"Edit the name of an existing TODO entry."`
	Delete  DeleteCmd  `cmd:"" aliases:"rm" help:"Delete existing TODO entries."`
	Archive ArchiveCmd `cmd:"" help:"Archive existing TODO entries."`
	Debug   bool       `help:"Enable debug mode for verbose logging."`
}

func PrintTodos(context Context, includeArchived bool) error {
	context.logger.LogDebug("retrieving all TODO entries")
	todos, err := context.db.GetAllTodoEntries()
	if err != nil {
		return err
	}

	context.logger.LogDebug(fmt.Sprintf("%d total TODO entries found", len(todos)))

	if includeArchived {
		fmt.Println("Active:")
	}

	for _, todo := range todos {
		if todo.status == Open && !todo.archived {
			fmt.Printf("%d: [ ] %s\n", todo.id, todo.name)
		}
	}

	for _, todo := range todos {
		if todo.status == Done && !todo.archived {
			fmt.Printf("%d: [X] %s\n", todo.id, todo.name)
		}
	}

	if includeArchived {
		fmt.Println("\nArchived:")

		for _, todo := range todos {
			if todo.archived {
				if todo.status == Open {
					fmt.Printf("%d: [ ] %s\n", todo.id, todo.name)
				} else {
					fmt.Printf("%d: [X] %s\n", todo.id, todo.name)
				}
			}
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
