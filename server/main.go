// Package main implements a server for the Wiz ProcessorAPI service.
package main

import (
	"github.com/alexkreidler/wiz/executor"
	"github.com/golang/protobuf/ptypes/empty"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	. "github.com/alexkreidler/wiz/api"
)

const (
	port = ":50051"
)

// server is used to implement the Server.
type server struct {
	UnimplementedProcessorAPIServer
	api executor.ProcessorExecutorAPI
}

func newServer() *server {
	return &server{
		api:executor.NewProcessorExecutor(),
	}
}


func (s server) GetAllProcessors(context.Context, *empty.Empty) (*Processors, error) {
	z:=s.api.GetAllProcessors()
	return &Processors{Processors:z}, nil
}

func (s server) GetProcessor(c context.Context, id *Processor_ID) (*Processor, error) {
	return s.api.GetProcessor(id.Id)
}

func (s server) GetRuns(context.Context, *Processor_ID) (*Run, error) {
	panic("implement me")
}

func (s server) GetRun(context.Context, *IndividualRunID) (*Configuration, error) {
	panic("implement me")
}

func (s server) GetConfig(context.Context, *IndividualRunID) (*Configuration, error) {
	panic("implement me")
}

func (s server) Configure(context.Context, *ConfigureRequest) (*empty.Empty, error) {
	panic("implement me")
}

func (s server) GetRunState(*IndividualRunID, ProcessorAPI_GetRunStateServer) error {
	panic("implement me")
}

func (s server) GetRunData(context.Context, *IndividualRunID) (*DataSpec, error) {
	panic("implement me")
}

func (s server) AddData(context.Context, *AddDataRequest) (*Data, error) {
	panic("implement me")
}


// SayHello implements helloworld.GreeterServer
// func (s *server) SayHello(ctx context.Context, in *pb.test) (*pb.HelloReply, error) {
// 	log.Printf("Received: %v", in.GetName())
// 	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
// }

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Listening on %s", port)
	s := grpc.NewServer()
	RegisterProcessorAPIServer(s, newServer())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}