package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	_ "modernc.org/sqlite"
)

var ValidCommands = []string{
	"list", "l",
	"new", "n",
	"delete", "d",
	"help", "h"}

func main() {
	if len(os.Args) < 2 {
		Help()
		os.Exit(0)
	}

	command := os.Args[1]
	if !slices.Contains(ValidCommands, command) {
		fmt.Printf("'%s' is not a valid command. See `todo help`\n", command)
		os.Exit(0)
	}

	db := GetOrCreateDatabase()
	defer db.Close()
	db.Init()

	switch strings.ToLower(command) {
	case "list", "l":
		List()
	case "new", "n":
		New()
	case "delete", "d":
		Delete()
	case "help", "h":
		Help()
	}
}

func List() {
	panic("TODO: Implement List")
}

func New() {
	panic("TODO: Implement New")
}

func Delete() {
	panic("TODO: Implement Delete")
}

func Help() {
	panic("TODO: Help Info")
}
