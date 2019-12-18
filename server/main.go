// Package main implements a server for the Wiz ProcessorAPI service.
package main

import (
	"context"
	"github.com/alexkreidler/wiz/executor"
	"github.com/gogo/protobuf/types"

	"log"
	"net"

	"github.com/alexkreidler/wiz/api"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// server is used to implement the Server.
type Server struct {
	//api.UnimplementedProcessorAPIServer
	api executor.ProcessorExecutorAPI
}

func (s Server) GetAllProcessors(context.Context, *types.Empty) (*api.Processors, error) {
	return &api.Processors{Processors: s.api.GetAllProcessors()}, nil
}

func (s Server) GetProcessor(c context.Context, id *api.ProcessorID) (*api.Processor, error) {
	p, err := s.api.GetProcessor(id.ID)
	if err != nil {
		return &api.Processor{}, err
	}
	return &p, nil
}

func (s Server) GetRuns(c context.Context, id *api.ProcessorID) (*api.Runs, error) {
	p, err := s.api.GetAllRuns(id.ID)
	if err != nil {
		return nil, err
	}
	return &api.Runs{Runs:p}, nil
}

func (s Server) GetRun(c context.Context, id *api.IndividualRunID) (*api.Run, error) {
	p, err := s.api.GetRun(id.ProcessorID.ID, id.ID)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s Server) GetConfig(c context.Context, id *api.IndividualRunID) (*api.Configuration, error) {
	p, err := s.api.GetConfig(id.ProcessorID.ID, id.ID)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s Server) Configure(c context.Context, req *api.ConfigureRequest) (*types.Empty, error) {
	_, err := s.api.UpdateConfig(req.RunID.ProcessorID.ID, req.RunID.ID, *req.Config)

	return &types.Empty{}, err
}

func (s Server) GetRunState(*api.IndividualRunID, api.ProcessorAPI_GetRunStateServer) error {
	//panic("implement me")
	log.Printf("not implemented Getrunstate")
	return nil
}

func (s Server) GetRunData(context.Context, *api.IndividualRunID) (*api.DataSpec, error) {
	//panic("implement me")
	log.Printf("not implemented GetrunData")
	return nil, nil
}

func (s Server) AddData(context.Context, *api.AddDataRequest) (*api.Data, error) {
	//panic("implement me")
	log.Printf("not implemented AddData")
	return nil, nil
}

func NewServer() *Server {
	return &Server{
		api:executor.NewProcessorExecutor(),
	}
}


func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Listening on %s", port)
	s := grpc.NewServer()

	srv := NewServer()
	api.RegisterProcessorAPIServer(s, srv)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}