package server

import (
	"fmt"
	"github.com/GearFramework/GophKeeper/internal/config"
	"github.com/caarlos0/env/v11"
)

const (
	defaultAddress       = ":8080"
	defaultLogLevel      = "info"
	defaultStorageDriver = ""
	defaultStorageDSN    = ""
	defaultMountPath     = "/tmp"
)

// NewConfig returning config of GopherKeeper server
func NewConfig(fl *Flags) (*config.Server, error) {
	if fl.ENVFile != empty {
		err := config.FromENV(fl.ENVFile)
		if err != nil {
			return nil, err
		}
	}
	conf := &config.Server{
		LogLevel:      defaultLogLevel,
		Addr:          defaultAddress,
		StorageDriver: defaultStorageDriver,
		StorageDSN:    defaultStorageDSN,
		MountPath:     defaultMountPath,
	}
	flagLoaded := false
	if fl.JSONFile != empty {
		if err := config.FromJSON(fl.JSONFile, conf); err != nil {
			return nil, err
		}
		flagLoaded = true
	}
	if fl.YAMLFile != empty {
		if !flagLoaded {
			return nil, config.ErrConflictFlags
		}
		if err := config.FromYAML(fl.YAMLFile, conf); err != nil {
			return nil, err
		}
	}
	mergeFlags(conf, fl)
	if err := env.Parse(conf); err != nil {
		return nil, err
	}
	fmt.Printf("Uses config: %v\n", conf)
	return conf, nil
}

func mergeFlags(conf *config.Server, fl *Flags) {
	if fl.Addr != empty {
		conf.Addr = fl.Addr
	}
	if fl.LogLevel != empty {
		conf.LogLevel = fl.LogLevel
	}
	if fl.StorageDriver != empty {
		conf.StorageDriver = fl.StorageDriver
	}
	if fl.StorageDSN != empty {
		conf.StorageDSN = fl.StorageDSN
	}
}
