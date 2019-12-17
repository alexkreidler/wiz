package executor

import (
	"fmt"
	"github.com/alexkreidler/wiz/api"
	"github.com/golang/protobuf/ptypes/any"
)

//ProcessorExecutorAPI implements the Wiz Processor interface, and can be called by an HTTP  or gRPC transport layer
type ProcessorExecutorAPI interface {
	GetAllProcessors() []*api.Processor
	GetProcessor(id string) (api.Processor, error)
	GetAllRuns(processorID string) ([]*api.Run, error)
	GetRun(processorID string, runID string) (api.Run, error)
	GetConfig(processorID string, runID string) (api.Configuration, error)
	UpdateConfig(processorID string, runID string, configuration api.Configuration) (api.Configuration, error)

	GetData(processorID string, runID string) (api.DataSpec, error)

	AddData(processorID string, runID string, chunkID string) (api.Data, error)
	GetInputChunk(processorID string, runID string, chunkID string) (api.Data, error)
	GetOutputChunk(processorID string, runID string, chunkID string) (api.Data, error)
}

const Version = "0.1.0"

//ProcessorExecutor implements the ProcessorExecutorAPI for builtin Golang Processors. It uses channels, maps, and concurrency to parallelize by chunks
type ProcessorExecutor struct {
	version string
	// p maps the ID of the processor to the processor. These are all base, non-configured processors that are registered at startup
	p map[string]Processor
	//concreteProcessors is a map of processor IDs to runIDs
	concreteProcessors map[string]map[string]Processor
}

func NewProcessorExecutor() *ProcessorExecutor {
	return &ProcessorExecutor{version: Version, p: make(map[string]Processor), concreteProcessors: make(map[string]map[string]Processor)}
}

//func NewProcessor

func (p ProcessorExecutor) GetAllProcessors() []*api.Processor {
	all := make([]*api.Processor, len(p.p))
	for _, processor := range p.p {
		z := processor.Metadata()
		all = append(all, &z)
	}
	return all
}

func (p ProcessorExecutor) GetProcessor(id string) (api.Processor, error) {
	processor, ok := p.p[id]
	if !ok {
		return api.Processor{}, fmt.Errorf("processor %s not found", id)
	} else {
		return processor.Metadata(), nil
	}
}

func (p ProcessorExecutor) GetAllRuns(processorID string) ([]*api.Run, error) {
	panic("implement me")
}

func (p ProcessorExecutor) GetRun(processorID string, runID string) (api.Run, error) {
	panic("implement me")
}

func (p ProcessorExecutor) GetConfig(processorID string, runID string) (api.Configuration, error) {
	_, ok := p.p[processorID]
	if !ok {
		return api.Configuration{}, fmt.Errorf("processor %s not found", processorID)
	}
	//return processor.Metadata()., nil
	return buildConfig(), nil
}
func buildConfig() api.Configuration  {
	return api.Configuration{Config: &any.Any{TypeUrl: "https://wiz-project.ml/configurtion",Value:[]byte(`test`)}}
}

func (p ProcessorExecutor) UpdateConfig(processorID string, runID string, configuration api.Configuration) (api.Configuration, error) {
	processor, ok := p.p[processorID]
	if !ok {
		return buildConfig(), fmt.Errorf("processor %s not found", processorID)
	}

	if p.concreteProcessors[processorID] == nil {
		return buildConfig(), fmt.Errorf("failed")
		//p.concreteProcessors[processorID] = make(map[string]Processor)
	}
	proc, err := processor.New(configuration)
	if err != nil {
		return buildConfig(), err
	}
	p.concreteProcessors[processorID][runID] = proc
	return configuration, nil
}

func (p ProcessorExecutor) GetData(processorID string, runID string) (api.DataSpec, error) {
	panic("implement me")
}

func (p ProcessorExecutor) AddData(processorID string, runID string, chunkID string) (api.Data, error) {
	panic("implement me")
}

func (p ProcessorExecutor) GetInputChunk(processorID string, runID string, chunkID string) (api.Data, error) {
	panic("implement me")
}

func (p ProcessorExecutor) GetOutputChunk(processorID string, runID string, chunkID string) (api.Data, error) {
	panic("implement me")
}
