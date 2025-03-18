package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	_ "modernc.org/sqlite"
)

type ListCmd struct {
}

func (l *ListCmd) Run(ctx *kong.Context) error {
	fmt.Printf("Inside the 'List' command\n")
	return nil
}

type NewCmd struct {
	Name string `short:"n" help:"The name of the TODO entry."`
}

func (n *NewCmd) Run(ctx *kong.Context) error {
	fmt.Printf("Inside the 'New' command\n")
	fmt.Printf("name=%s\n", n.Name)
	return nil
}

type DeleteCmd struct {
}

func (d *DeleteCmd) Run(ctx *kong.Context) error {
	fmt.Printf("Inside the 'Delete' command\n")
	return nil
}

var CLI struct {
	List   ListCmd   `cmd:"" help:"List TODO entries."`
	New    NewCmd    `cmd:"" help:"Create a new TODO entry."`
	Delete DeleteCmd `cmd:"" help:"Delete an existing TODO entry."`
}

func main() {
	db := GetOrCreateDatabase()
	defer db.Close()
	db.Init()

	ctx := kong.Parse(&CLI)
	ctx.Run()
}
