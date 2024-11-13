package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/GearFramework/GophKeeper/internal/client"
	"log"
	"os"
)

func printUsage() {
	fmt.Print(`
Usage: 
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

var (
	BuildVersion string
	BuildDate    string
)

func stringBuild(b string) string {
	if b == "" {
		return "N/A"
	}
	return b
}

func printGreeting() {
	fmt.Printf("Build version: %s\nBuild date: %s\n",
		stringBuild(BuildVersion),
		stringBuild(BuildDate),
	)
}

func main() {
	printGreeting()
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
	tlsConf, err := getTLSConfig()
	if err != nil {
		return err
	}
	c := client.GkClient{
		Conf: fl,
		Tls:  tlsConf,
	}
	if err != nil {
		return err
	}
	return c.Run()
}

var tlsCertFile = ".cert/certbundle.pem"

func getTLSConfig() (*tls.Config, error) {
	//http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	caCert, err := os.ReadFile(tlsCertFile)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	return &tls.Config{
		RootCAs: caCertPool,
	}, nil
}
