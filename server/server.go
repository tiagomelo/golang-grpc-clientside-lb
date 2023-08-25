// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.
//
// Package server provides the implementation of the gRPC server functionalities
// for the Hello Service, including starting the server and handling RPCs.
package server

import (
	"context"
	"fmt"

	"github.com/tiagomelo/golang-grpc-clientside-lb/api/proto/gen/helloservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// server struct encapsulates the gRPC server functionalities and
// implements the GreeterServer interface from the proto definition.
type server struct {
	helloservice.UnimplementedGreeterServer
	GrpcSrv *grpc.Server
	host    string
}

// New initializes and returns a new instance of the server with the provided host string.
func New(host string) *server {
	grpcServer := grpc.NewServer()
	srv := &server{
		GrpcSrv: grpcServer,
		host:    host,
	}
	helloservice.RegisterGreeterServer(grpcServer, srv)
	reflection.Register(grpcServer)
	return srv
}

// SayHello is an RPC method implementation that responds with a greeting message,
// including the server's host identifier.
func (s *server) SayHello(ctx context.Context, in *helloservice.HelloRequest) (*helloservice.HelloResponse, error) {
	return &helloservice.HelloResponse{Message: fmt.Sprintf("Hello, %s (from server %s)", in.Name, s.host)}, nil
}
