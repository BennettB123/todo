package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

func CreateDirectory(path string) {
	_, err := os.ReadDir(path)
	if errors.Is(err, fs.ErrNotExist) {
		fmt.Printf("creating directory at '%s'\n", path)
		err = os.Mkdir(path, 0700)
		if err != nil {
			fmt.Printf("unable to create directory at '%s'.\n", path)
			fmt.Printf("  error: %s\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("directory at '%s' already exists\n", path)
	}
}