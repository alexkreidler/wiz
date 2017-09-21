package daemon

import (
	"log"
	"net"

	pb "github.com/tim15/wiz/api/daemon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

func (s *server) GetVersion(ctx context.Context, in *pb.Empty) (*pb.Version, error) {
	return &pb.Version{Version: "0.0.1"}, nil
}

func (s *server) InstallPackages(ctx context.Context, in *pb.PackageList) {
	log.Printf("Installing packages %+v\n", in)
}

func (s *server) GetPackages(ctx context.Context, in *pb.Empty) {
	log.Println("Getting packages")
	return pb.PackageList{}
}

func Start() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDaemonServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
