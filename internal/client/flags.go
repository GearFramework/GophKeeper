package client

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"slices"
)

// Command cli-command in flag -c
type Command string

const (
	// CommandSignup register new customer
	CommandSignup Command = "signup"
	// CommandAdd create new entity into remote server
	CommandAdd Command = "add"
	// CommandDel delete entity on remote server
	CommandDel Command = "del"
	// CommandList get stored entities on remote server
	CommandList Command = "list"
	// CommandView get entity and display metadata
	CommandView Command = "view"
	// CommandUpload upload binary file of exists entity
	CommandUpload Command = "upload"
	// CommandDownload download binary file of entity
	CommandDownload Command = "download"
)

// NeedAuth return true if command need authorization on remote server
func (c Command) NeedAuth() bool {
	return c != CommandSignup
}

var (
	commmands = []Command{
		CommandSignup,
		CommandAdd,
		CommandDel,
		CommandList,
		CommandView,
		CommandUpload,
		CommandDownload,
	}
)

var (
	// ErrInvalidCommand raised if unknown command in flag -c
	ErrInvalidCommand = errors.New("invalid command")
)

// ParseFlags parse command-line flags
func ParseFlags() (*Config, error) {
	fl := Config{}
	if len(os.Args) < 2 {
		return nil, ErrInvalidCommand
	}
	//command := Command(os.Args[1][1:])
	var cmd string
	fmt.Printf("Usage of %s:\n", fl.Command)
	fmt.Printf("Service started with flags: %v\n", os.Args[1:])
	flag.StringVar(&fl.Addr, "a", "", "address to remote server")
	flag.StringVar(&fl.Username, "u", "", "authentication username")
	flag.StringVar(&cmd, "c", "", "command to execute")
	flag.StringVar(&fl.Type, "t", "", "data type")
	flag.Parse()
	fl.Command = Command(cmd)
	if !slices.Contains(commmands, fl.Command) {
		return nil, ErrInvalidCommand
	}
	fmt.Println("Config from flags: ", fl)
	return &fl, nil
}
