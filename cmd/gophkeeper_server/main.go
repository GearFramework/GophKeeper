package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GearFramework/GophKeeper/internal/pkg/logger"
	"github.com/GearFramework/GophKeeper/internal/server"
)

var (
	buildVersion string = "0.0.1"
	buildDate    string = time.Date(2024, 11, 10, 16, 14, 0, 0, time.UTC).Format(time.RFC3339)
	buildCommit  string = "1"
)

func stringBuild(b string) string {
	if b == "" {
		return "N/A"
	}
	return b
}

func printGreeting() {
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n",
		stringBuild(buildVersion),
		stringBuild(buildDate),
		stringBuild(buildCommit),
	)
}

func main() {
	printGreeting()
	if err := run(); err != nil {
		log.Fatal(err.Error())
	}
}

func run() error {
	fl := server.ParseFlags()
	conf, err := server.NewConfig(fl)
	if err != nil {
		return err
	}
	if err = logger.Init(conf.LogLevel); err != nil {
		return err
	}
	app := server.NewGkServer(conf)
	if err = app.Init(); err != nil {
		logger.Log.Errorf("Error initializing server: %s", err.Error())
		return err
	}
	gracefulStop(func() {
		app.Stop()
	})
	return app.Run()
}

func gracefulStop(stopCallback func()) {
	gracefulStopChan := make(chan os.Signal, 1)
	signal.Notify(
		gracefulStopChan,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)
	go func() {
		sig := <-gracefulStopChan
		stopCallback()
		log.Printf("Caught sig: %+v\n", sig)
		log.Println("Application graceful stop!")
		os.Exit(0)
	}()
}
