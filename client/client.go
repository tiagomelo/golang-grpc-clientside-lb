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
func makeRPCs(logger *log.Logger, cc *grpc.ClientConn, n int) error {
	client := helloservice.NewGreeterClient(cc)
	for i := 0; i < n; i++ {
		message, err := sayHello(client, "Tiago")
		if err != nil {
			return errors.Wrap(err, "calling SayHello")
		}
		fmt.Println(message)
		time.Sleep(1 * time.Second)
	}
	return nil
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
	return makeRPCs(logger, roundrobinConn, 10)
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
