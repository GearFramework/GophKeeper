package config

import (
	"errors"

	"github.com/GearFramework/GophKeeper/internal/gk"
	"github.com/joho/godotenv"
)

var (
	// ErrEnvFileNotFound file .env not found by path
	ErrEnvFileNotFound = errors.New("env file not found")
)

// FromENV set environment from .env file
func FromENV(filepath string) error {
	if !gk.IsExistsFile(filepath) {
		return ErrEnvFileNotFound
	}
	return godotenv.Overload(filepath)
}
