// Package main implements a server for the Wiz ProcessorAPI service.
package main

import (
	"github.com/alexkreidler/wiz/executor"
	"log"
	"net"

	"github.com/alexkreidler/wiz/api"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)


func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Listening on %s", port)
	s := grpc.NewServer()

	srv := executor.NewProcessorExecutor()
	api.RegisterProcessorAPIServer(s, srv)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}