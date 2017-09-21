package daemon

import (
	"log"
	"net"

	"context"
	"github.com/tim15/wiz/api/daemon"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

func (s *server) GetVersion(ctx context.Context, in *daemon.Empty) (*daemon.Version, error) {
	return &daemon.Version{Version: "0.0.1"}, nil
}

func (s *server) InstallPackages(ctx context.Context, in *daemon.PackageList) (*daemon.Status, error) {
	log.Printf("Installing packages %+v\n", in)
	return InstallPackages(*in), nil
}

func (s *server) GetPackages(ctx context.Context, in *daemon.Empty) (*daemon.PackageList, error) {
	log.Println("Getting packages")
	return &daemon.PackageList{}, nil
}

func (s *server) GetConfig(ctx context.Context, in *daemon.Empty) (*daemon.Config, error) {
	return &daemon.Config{PackageLocation: "/tmp/wiz/packages", UseComputeDevices: "all", RunBackend: daemon.Config_LOCAL}, nil
}

func (s *server) SetConfig(context.Context, *daemon.Config) (*daemon.Status, error) {
	return &daemon.Status{Status: true}, nil
}

func Start() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	daemon.RegisterDaemonServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
