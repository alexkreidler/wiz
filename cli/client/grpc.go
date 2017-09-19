package client

import (
	pb "github.com/tim15/wiz/api/daemon"
	"google.golang.org/grpc"
	"log"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func GetClient() (pb.DaemonClient, error) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return nil, err
	}
	c := pb.NewDaemonClient(conn)

	return c, nil
}
