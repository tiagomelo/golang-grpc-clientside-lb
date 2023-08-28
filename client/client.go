// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.
//
// Package main provides a gRPC client implementation that establishes a connection to a
// server and initiates RPC calls.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
	"github.com/tiagomelo/golang-grpc-clientside-lb/api/proto/gen/helloservice"
	"github.com/tiagomelo/golang-grpc-clientside-lb/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// sayHello sends a greeting request with the given name and retrieves the server's response.
func sayHello(c helloservice.GreeterClient, name string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &helloservice.HelloRequest{Name: name})
	if err != nil {
		return "", err
	}
	return r.Message, nil
}

// makeRPCs sends a series of RPC calls to the server and prints the response.
func makeRPCs(logger *log.Logger, cc *grpc.ClientConn) error {
	client := helloservice.NewGreeterClient(cc)
	for {
		message, err := sayHello(client, "Tiago")
		if err != nil {
			return errors.Wrap(err, "calling SayHello")
		}
		fmt.Println(message)
		time.Sleep(1 * time.Second)
	}
}

// run sets up the gRPC client, establishes a connection to the server, and initiates RPC calls.
func run(logger *log.Logger, loadBalancingPolicy string) error {
	defer logger.Println("main: Completed")
	cfg, err := config.Read()
	if err != nil {
		return errors.Wrap(err, "reading config")
	}
	grpcServiceConfig := fmt.Sprintf(`{"loadBalancingConfig": [{"%s":{}}]}`, loadBalancingPolicy)
	roundrobinConn, err := grpc.Dial(
		cfg.ServiceTargetAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(grpcServiceConfig),
	)
	if err != nil {
		return errors.Wrap(err, "dialing")
	}
	defer roundrobinConn.Close()

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Send requests.
	go func() {
		serverErrors <- makeRPCs(logger, roundrobinConn)
	}()
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		logger.Println("main: received signal for shutdown: ", sig)
	}

	return nil
}

// options struct holds command line flags configurations.
type options struct {
	LoadBalancingPolicy string `short:"l" description:"load balancing policy" choice:"round_robin" choice:"pick_first"`
}

func main() {
	var opts options
	flags.Parse(&opts)
	logger := log.New(os.Stdout, "HELLO SERVICE CLIENT : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	if err := run(logger, opts.LoadBalancingPolicy); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
