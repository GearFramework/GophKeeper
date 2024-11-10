package config

import (
	"errors"
)

var (
	// ErrInvalidConfigFile raise if config file not found
	ErrInvalidConfigFile = errors.New("invalid config file")
	// ErrConflictFlags raise if
	ErrConflictFlags = errors.New("conflict flags -j and -y of file configs")
)

// Server struct of config
type Server struct {
	LogLevel      string `json:"log_level",yaml:"log_level",env:"GK_LOG_LEVEL,required"`
	Addr          string `json:"addr",yaml:"addr",env:"GK_SERVER_ADDR,required"`
	StorageDriver string `json:"storage_driver",yaml:"storage_driver",env:"GK_SERVER_STORAGE_DRIVER,required"`
	StorageDSN    string `json:"storage_dsn",yaml:"storage_dsn",env:"GK_SERVER_STORAGE_DSN,required"`
	MountPath     string `json:"mount_path",yaml:"mount_path",env:"GK_SERVER_MOUNT_PATH,required"`
}
