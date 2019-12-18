package executor

import (
	"context"
	"fmt"
	"github.com/alexkreidler/wiz/api"
	"github.com/alexkreidler/wiz/processors"
	"github.com/gogo/protobuf/types"
)

const Version = "0.1.0"

//runProcessor contains a run and a processor which is that run
type runProcessor struct {
	// r is the metadata about the run for serialization
	r api.Run

	// p is the actual instance of the processor that has been configured
	p processors.Processor

	// ds holds all of the processor's data
	ds api.DataSpec
}

//ProcessorExecutor implements the ProcessorAPIServer for builtin Golang Processors. It uses channels, maps, and concurrency to parallelize by chunks
type ProcessorExecutor struct {
	version string
	// base maps the ID of the processor to the processor. These are all base, non-configured processors that are registered at startup
	//Their Metadata functions are the source of all processor information
	base map[string]processors.Processor
	//runMap is a map of processor IDs to runIDs and processors
	runMap runProcMap
}

type runProcMap map[string]map[string]*runProcessor

func (p ProcessorExecutor) GetAllProcessors(context.Context, *types.Empty) (*api.Processors, error) {
	all := make([]api.Processor, len(p.base))
	for _, processor := range p.base {
		all = append(all, processor.Metadata())
	}
	return &api.Processors{all}, nil
}

// TODO: update to use new data access
func (p ProcessorExecutor) GetProcessor(c context.Context, id *api.ProcessorID) (*api.Processor, error) {
	processor, ok := p.base[id.ID]
	if !ok {
		return nil, fmt.Errorf("processor %s not found", id.ID)
	} else {
		p := processor.Metadata()
		return &p, nil
	}
}

func (p ProcessorExecutor) GetRuns(c context.Context, id *api.ProcessorID) (*api.Runs, error) {
	processor, ok := p.runMap[id.ID]
	if !ok {
		return nil, fmt.Errorf("either processor %s does not exist or no runs are registered", id.ID)
	} else {
		all := make([]api.Run, len(processor))
		for _, run := range processor {
			all = append(all, run.r)
		}
		return &api.Runs{Runs: all}, nil
	}
}

func (p ProcessorExecutor) GetRun(c context.Context, req *api.IndividualRunID) (*api.Run, error) {
	r, err := getRun(p, req.ID, req.RunID)
	if err != nil {
		return nil, err
	}
	return &r.r, nil
}

// Data access utilities
func checkProcessorExists(p ProcessorExecutor, processorID string) error {
	_, ok := p.base[processorID]
	if !ok {
		return fmt.Errorf("processor %s not found", processorID)
	}
	return nil
}

func chuckRunExists(p ProcessorExecutor, id string, runID string) error {
	if err := checkProcessorExists(p, id); err != nil {
		return err
	}
	_, ok := p.runMap[id][runID]
	if !ok {
		return fmt.Errorf("run %s not found", runID)
	}
	return nil
}

//getRun returns a pointer because the run, processor, and data need to be updated
func getRun(p ProcessorExecutor, id string, runID string) (*runProcessor, error) {
	err := chuckRunExists(p, id, runID)
	if err != nil {
		return nil, err
	}

	return p.runMap[id][runID], nil
}

func (p ProcessorExecutor) GetConfig(c context.Context, id *api.IndividualRunID) (*api.Configuration, error) {
	r, err := getRun(p, id.ID, id.RunID)
	if err != nil {
		return nil, err
	}
	return r.r.Config, nil
}

func (p ProcessorExecutor) Configure(c context.Context, req *api.ConfigureRequest) (*types.Empty, error) {
	processorID := req.RunID.ID
	runID := req.RunID.RunID
	configuration := req.Config.Config
	//TODO: unmarshal config

	baseProcessor, ok := p.base[processorID]
	if !ok {
		return nil, fmt.Errorf("baseProcessor %s not found", processorID)
	}

	if p.runMap[processorID] == nil {
		return nil, fmt.Errorf("failed")
		//base.concreteProcessors[processorID] = make(map[string]Processor)
	}
	processor, err := baseProcessor.New(configuration)
	if err != nil {
		return nil, err
	}
	p.runMap[processorID][runID] = &runProcessor{p: processor, r: api.Run{
		RunID: runID,
		//TODO: think about storing the config in a deserialized format here for easy access. Then again, the processor's already configured and shouldn't need to be again
		Config: req.Config,
		State:  api.State_CONFIGURED,
	}}
	return nil, nil
}

func (p ProcessorExecutor) GetRunState(*api.IndividualRunID, api.ProcessorAPI_GetRunStateServer) error {
	panic("implement me")
}

func (p ProcessorExecutor) GetRunData(c context.Context, id *api.IndividualRunID) (*api.DataSpec, error) {
	r, err := getRun(p, id.ID, id.RunID)
	if err != nil {
		return nil, err
	}
	return &r.ds, nil
}

func (p ProcessorExecutor) AddData(c context.Context, req *api.AddDataRequest) (*api.Data, error) {
	r, err := getRun(p, req.Id.ID, req.Id.RunID)
	if err != nil {
		return nil, err
	}
	//todo: make all these chunks concurrent
	r.ds.In = append(r.ds.In, req.Data)
	return nil, nil
}

func NewProcessorExecutor() ProcessorExecutor {
	return ProcessorExecutor{version: Version, base: processors.DefaultProcessors, runMap: make(runProcMap)}
}
//
//func buildConfig() *api.Configuration {
//	return &api.Configuration{Config: &types.Any{TypeUrl: "https://wiz-project.ml/configurtion", Value: []byte(`test`)}}
//}
