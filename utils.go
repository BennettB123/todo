package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

func CreateDirectory(path string, logger Logger) error {
	logger.LogDebug(fmt.Sprintf("checking if directory exists: %s", path))
	_, err := os.ReadDir(path)
	if errors.Is(err, fs.ErrNotExist) {
		logger.LogDebug("directory does not exist; creating...")
		err = os.Mkdir(path, 0700)
		if err == nil {
			logger.LogDebug("directory creation successful")
		}
		return err
	} else {
		logger.LogDebug("directory already exists")
	}

	return nil
}
