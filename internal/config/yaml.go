package config

import (
	"os"

	"github.com/GearFramework/GophKeeper/internal/gk"
	"gopkg.in/yaml.v3"
)

// FromYAML get config from YAML-file
func FromYAML(filepath string, v interface{}) error {
	if !gk.IsExistsFile(filepath) {
		return ErrInvalidConfigFile
	}
	b, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(b, v); err != nil {
		return err
	}
	return nil
}
