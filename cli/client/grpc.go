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
	defer conn.Close()
	c := pb.NewDaemonClient(conn)

	return c, nil
	// Contact the server and print out its response.
	// name := defaultName
	// if len(os.Args) > 1 {
	//   name = os.Args[1]
	// }
	// r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
	// if err != nil {
	//   log.Fatalf("could not greet: %v", err)
	// }
	// log.Printf("Greeting: %s", r.Message)
}
