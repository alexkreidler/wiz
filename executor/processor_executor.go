package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alexkreidler/wiz/api"
	"github.com/alexkreidler/wiz/processors"
	procApi "github.com/alexkreidler/wiz/processors/processor"
	"github.com/davecgh/go-spew/spew"
	"github.com/gogo/protobuf/types"
	"log"
)

const Version = "0.1.0"

//runProcessor contains a run and a processor which is that run
type runProcessor struct {
	// r is the metadata about the run for serialization
	r api.Run

	// p is the actual instance of the processor that has been configured
	p procApi.Processor

	// ds holds all of the processor's data
	ds api.DataSpec
}

//ProcessorExecutor implements the ProcessorAPIServer for builtin Golang Processors. It uses channels, maps, and concurrency to parallelize by chunks
type ProcessorExecutor struct {
	version string
	// base maps the ID of the processor to the processor. These are all base, non-configured processors that are registered at startup
	//Their Metadata functions are the source of all processor information
	base map[string]procApi.Processor
	//runMap is a map of processor IDs to runIDs and processors
	runMap runProcMap
}

type runProcMap map[string]map[string]*runProcessor

func (p ProcessorExecutor) GetAllProcessors(context.Context, *types.Empty) (*api.Processors, error) {
	all := make([]api.Processor, 0)
	for _, processor := range p.base {
		all = append(all, processor.Metadata())
	}
	return &api.Processors{all}, nil
}

// TODO: update to use new data access
func (p ProcessorExecutor) GetProcessor(c context.Context, id *api.ProcessorID) (*api.Processor, error) {
	processor, ok := p.base[id.ProcID]
	if !ok {
		return nil, fmt.Errorf("processor %s not found", id.ProcID)
	} else {
		p := processor.Metadata()
		return &p, nil
	}
}

func (p ProcessorExecutor) GetRuns(c context.Context, id *api.ProcessorID) (*api.Runs, error) {
	err := checkProcessorExists(p, id.ProcID)
	if err != nil {
		return nil, err
	}
	processor, ok := p.runMap[id.ProcID]
	if !ok {
		// no runs are registered
		return &api.Runs{}, nil
	} else {
		all := make([]api.Run, len(processor))
		for _, run := range processor {
			all = append(all, run.r)
		}
		return &api.Runs{Runs: all}, nil
	}
}

func (p ProcessorExecutor) GetRun(c context.Context, req *api.IndividualRunID) (*api.Run, error) {
	r, err := getRun(p, req.ProcID, req.RunID)
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
	r, err := getRun(p, id.ProcID, id.RunID)
	if err != nil {
		return nil, err
	}
	return &api.Configuration{r.r.Config}, nil
}

func (p ProcessorExecutor) Configure(c context.Context, req *api.ConfigureRequest) (*types.Empty, error) {
	runID := req.ID.RunID
	processorID := req.ID.ProcID
	configuration := req.Config
	//TODO: unmarshal config

	baseProcessor, ok := p.base[processorID]
	if !ok {
		return nil, fmt.Errorf("baseProcessor %s not found", processorID)
	}

	if p.runMap[processorID] == nil {
		//return nil, fmt.Errorf("failed")
		log.Printf("proc %s did not have any runs, creating", processorID)
		p.runMap[processorID] = make(map[string]*runProcessor)
	}

	spew.Dump(configuration)

	//The run won't exist here at this point, so we create it:
	proc, err := baseProcessor.New(configuration)
	if err != nil {
		return nil, err
	}
	fmt.Println("got here")
	rp := &runProcessor{p: proc, r: api.Run{
		RunID: runID,
		//TODO: think about storing the config in a deserialized format here for easy access. Then again, the proc's already configured and shouldn't need to be again
		Config: req.Config,
		State:  api.State_CONFIGURED,
	}}

	spew.Dump(rp)

	p.runMap[processorID][runID] = rp
	return &types.Empty{}, nil
}

func (p ProcessorExecutor) GetRunState(*api.IndividualRunID, api.ProcessorAPI_GetRunStateServer) error {
	panic("implement me")
}

func (p ProcessorExecutor) GetRunData(c context.Context, id *api.IndividualRunID) (*api.DataSpec, error) {
	r, err := getRun(p, id.ProcID, id.RunID)
	if err != nil {
		return nil, err
	}
	return &r.ds, nil
}

func (p ProcessorExecutor) AddData(c context.Context, req *api.AddDataRequest) (*types.Empty, error) {
	r, err := getRun(p, req.ID.ProcID, req.ID.RunID)
	if err != nil {
		return nil, err
	}
	//todo: make all these chunks concurrent
	r.ds.In = append(r.ds.In, req.Data)
	//switch req.Data.Data.
	d, ok := req.Data.Data.(*api.Data_Raw)
	if !ok {
		return nil, fmt.Errorf("failed to read raw data")
	}
	val := (*d).Raw.Value
	var data interface{}
	spew.Dump(val)
	err = json.Unmarshal(val, &data)
	if err != nil {
		return nil, err
	}
	//todo: unmarshall json
	r.p.Run(data)
	return nil, nil
}

func NewProcessorExecutor() ProcessorExecutor {
	initProc := processors.ConfiguredProcessorRegistry().Processors
	spew.Dump(initProc)
	return ProcessorExecutor{version: Version, base: initProc, runMap: make(runProcMap)}
}
//
//func buildConfig() *api.Configuration {
//	return &api.Configuration{Config: &types.Any{TypeUrl: "https://wiz-project.ml/configurtion", Value: []byte(`test`)}}
//}
