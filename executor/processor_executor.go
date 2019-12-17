package executor

import (
	"fmt"
	"github.com/alexkreidler/wiz/models"
)

//ProcessorExecutorAPI implements the Wiz Processor interface, and can be called by an HTTP  or gRPC transport layer
type ProcessorExecutorAPI interface {
	GetAllProcessors() []*models.ProcessorObject
	GetProcessor(id string) (models.ProcessorObject, error)
	GetAllRuns(processorID string) ([]*models.Run, error)
	GetRun(processorID string, runID string) (models.Run, error)
	GetConfig(processorID string, runID string) (models.Configuration, error)
	UpdateConfig(processorID string, runID string, configuration models.Configuration) (models.Configuration, error)

	GetData(processorID string, runID string) (models.DataSpec, error)

	AddData(processorID string, runID string, chunkID string) (models.Data, error)
	GetInputChunk(processorID string, runID string, chunkID string) (models.Data, error)
	GetOutputChunk(processorID string, runID string, chunkID string) (models.Data, error)
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

func (p ProcessorExecutor) GetAllProcessors() []*models.ProcessorObject {
	all := make([]*models.ProcessorObject, len(p.p))
	for _, processor := range p.p {
		z := processor.Metadata()
		all = append(all, &z)
	}
	return all
}

func (p ProcessorExecutor) GetProcessor(id string) (models.ProcessorObject, error) {
	processor, ok := p.p[id]
	if !ok {
		return models.ProcessorObject{}, fmt.Errorf("processor %s not found", id)
	} else {
		return processor.Metadata(), nil
	}
}

func (p ProcessorExecutor) GetAllRuns(processorID string) ([]*models.Run, error) {
	panic("implement me")
}

func (p ProcessorExecutor) GetRun(processorID string, runID string) (models.Run, error) {
	panic("implement me")
}

func (p ProcessorExecutor) GetConfig(processorID string, runID string) (models.Configuration, error) {
	processor, ok := p.p[processorID]
	if !ok {
		return nil, fmt.Errorf("processor %s not found", processorID)
	}
	return processor.Metadata(), nil
}

func (p ProcessorExecutor) UpdateConfig(processorID string, runID string, configuration models.Configuration) (models.Configuration, error) {
	processor, ok := p.p[processorID]
	if !ok {
		return nil, fmt.Errorf("processor %s not found", processorID)
	}

	if p.concreteProcessors[processorID] == nil {
		return nil, fmt.Errorf("failed")
		//p.concreteProcessors[processorID] = make(map[string]Processor)
	}
	proc, err := processor.New(configuration)
	if err != nil {
		return nil, err
	}
	p.concreteProcessors[processorID][runID] = proc
	return configuration, nil
}

func (p ProcessorExecutor) GetData(processorID string, runID string) (models.DataSpec, error) {
	panic("implement me")
}

func (p ProcessorExecutor) AddData(processorID string, runID string, chunkID string) (models.Data, error) {
	panic("implement me")
}

func (p ProcessorExecutor) GetInputChunk(processorID string, runID string, chunkID string) (models.Data, error) {
	panic("implement me")
}

func (p ProcessorExecutor) GetOutputChunk(processorID string, runID string, chunkID string) (models.Data, error) {
	panic("implement me")
}
