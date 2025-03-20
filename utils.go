package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

func CreateDirectory(path string, logger Logger) error {
	logger.Log(Debug, fmt.Sprintf("checking if directory exists: %s", path))
	_, err := os.ReadDir(path)
	if errors.Is(err, fs.ErrNotExist) {
		logger.Log(Debug, "directory does not exist; creating...")
		err = os.Mkdir(path, 0700)
		if err == nil {
			logger.Log(Debug, "directory creation successful")
		}
		return err
	} else {
		logger.Log(Debug, "directory already exists")
	}

	return nil
}
