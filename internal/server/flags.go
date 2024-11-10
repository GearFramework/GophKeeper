package server

import (
	"flag"
	"fmt"
	"os"

	"github.com/GearFramework/GophKeeper/internal/config"
)

var empty = "\x00"

// Flags struct of command-line flags
type Flags struct {
	config.Server
	ENVFile  string
	JSONFile string
	YAMLFile string
}

// ParseFlags parse command-line flags
func ParseFlags() *Flags {
	var fl Flags
	fmt.Printf("Service started with flags: %v\n", os.Args[1:])
	flag.StringVar(&fl.Addr, "a", empty, "address to run server")
	flag.StringVar(&fl.LogLevel, "l", empty, "logger level")
	flag.StringVar(&fl.StorageDriver, "s", empty, "storage driver name")
	flag.StringVar(&fl.StorageDSN, "d", empty, "storage connection DSN")
	flag.StringVar(&fl.ENVFile, "e", empty, ".env config file")
	flag.StringVar(&fl.JSONFile, "j", empty, ".json config file")
	flag.StringVar(&fl.YAMLFile, "y", empty, ".yaml config file")
	flag.Parse()
	fmt.Println("Config from flags: ", fl)
	return &fl
}
