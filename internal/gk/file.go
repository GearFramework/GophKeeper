package gk

import (
	"errors"
	"os"
)

// IsExistsFile return true if file exists else returning false
func IsExistsFile(filepath string) bool {
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
