package main

import (
	"github.com/jessjenkins/github-bots/service"
	"github.com/pkg/errors"
	"log"
	"os"
	"os/signal"
)

func main() {

	if err := run(); err != nil {
		log.Fatalf("fatal runtime error: %v\n", err)
	}
}

func run() error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)

	svc, err := service.Create()
	if err != nil {
		return errors.Wrap(err, "creation of service failed")
	}

	svcErrors := make(chan error, 1)
	err = svc.Run(svcErrors)
	if err != nil {
		return errors.Wrap(err, "running service failed")
	}

	// blocks until an os interrupt or a fatal error occurs
	select {
	case err := <-svcErrors:
		return errors.Wrap(err, "service error received")
	case sig := <-signals:
		log.Printf("os signal received: %v", sig)
		svc.Close()
	}
	return nil
}
