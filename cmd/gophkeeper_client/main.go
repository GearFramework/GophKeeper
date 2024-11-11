package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/GearFramework/GophKeeper/internal/client"
)

func printUsage() {
	fmt.Println(`Usage: 
	gophkeeper [arguments]

The arguments are:

	-c <command>, possible commands
		add 		create new entity and send to remote server
		del		delete entity on remote server
		list		list all entities on remote server
		view		get entity from remote server and show meta data
		upload		upload file to remote server
		download	get file from remote server
	-a <address> addres of remote server, example: 
		http://localhost:8080
	-t <type of new entity> possible types:
		text		multi-line text data
		credentials	login and password data
		creditcard	data of credit cards
		binary		various binary file
	-u <username>
`)
}

var (
	errInvalidFlags = errors.New("invalid flags")
)

func main() {
	if err := run(); err != nil {
		if errors.Is(err, errInvalidFlags) {
			printUsage()
		} else {
			log.Fatal(err.Error())
		}
	}
}

func run() error {
	fl, err := client.ParseFlags()
	if err != nil {
		return errInvalidFlags
	}
	c := client.GkClient{
		Conf: fl,
	}
	if err != nil {
		return err
	}
	return c.Run()
}
