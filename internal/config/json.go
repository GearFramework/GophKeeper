package config

import (
	"encoding/json"
	"os"

	"github.com/GearFramework/GophKeeper/internal/gk"
)

// FromJSON get config from JSON-file
func FromJSON(filepath string, v interface{}) error {
	if !gk.IsExistsFile(filepath) {
		return ErrInvalidConfigFile
	}
	b, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, v); err != nil {
		return err
	}
	return nil
}
