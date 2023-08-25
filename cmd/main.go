// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.
package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
	"github.com/tiagomelo/golang-grpc-clientside-lb/server"
)

func run(logger *log.Logger, serverHost string, serverPort int) error {
	defer logger.Println("main: Completed")

	// Setting up a TCP listener on the specified port.
	port := fmt.Sprintf("%s:%d", serverHost, serverPort)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return errors.Wrap(err, "tcp listening")
	}

	// Initialize the gRPC server.
	srv := server.New(port)

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.

	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		logger.Printf("main: gRPC server listening on %s", port)
		serverErrors <- srv.GrpcSrv.Serve(lis)
	}()
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		logger.Println("main: received signal for shutdown: ", sig)
		srv.GrpcSrv.Stop()
	}

	return nil
}

// options struct holds command line flags configurations.
type options struct {
	ServerHost string `short:"s" description:"server host"`
	ServerPort int    `short:"p" description:"server port"`
}

func main() {
	var opts options
	flags.Parse(&opts)
	logger := log.New(os.Stdout, "HELLO SERVICE SERVER : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	if err := run(logger, opts.ServerHost, opts.ServerPort); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
