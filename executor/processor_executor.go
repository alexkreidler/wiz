package executor

import (
	"fmt"
	"github.com/alexkreidler/wiz/api"
	"github.com/alexkreidler/wiz/processors"
	procApi "github.com/alexkreidler/wiz/processors/processor"
	"github.com/davecgh/go-spew/spew"
	"log"
	"sync"
)

const Version = "0.1.0"

//runProcessor contains a run and a processor which is that run
type runProcessor struct {
	baseProcessor procApi.ChunkProcessor

	// runLock locks the run information (state specifically).
	// TODO: think about locking configuration. currently we don't allow configuration changes
	runLock sync.RWMutex
	// run is the metadata about the run for serialization
	run api.Run

	// dataLock locks the state of all the data chunks/workers
	dataLock sync.RWMutex
	// workers is a map from ChunkID to worker
	workers map[string]*Worker


	// a WaitGroup doesn't make sense cause it needs the processor goroutine which we don't control to cal  Done. We could just wrap it?? in an anon goroutine

	numCompleted uint32
}

type Worker struct {
	p   procApi.ChunkProcessor
	in  api.Data
	out api.Data
}

type runProcMap map[string]map[string]*runProcessor

//ProcessorExecutor implements the ProcessorAPIServer for builtin Golang Processors. It uses channels, maps, and concurrency to parallelize by chunks
//It is designed so all operations can use a value receiver. The version an base processors don't need to be modified, but the specific runs do, which is why it is a map to pointers
//can maps be modified with value receivers? yes they can, because maps, slices, and channels are inherently mutable
//think about concurrency issues, e.g. will multiple requests result in consistent state -- map need a stronger concurrency primitive than a simple Map, e.g. locks
type ProcessorExecutor struct {
	version string
	// base maps the ID of the processor to the processor. These are all base, non-configured processors that are registered at startup
	//Their Metadata functions are the source of all processor information
	base map[string]procApi.ChunkProcessor
	//runMap is a map of processor IDs to runIDs and processors
	runMap runProcMap
}

func (p ProcessorExecutor) GetAllProcessors() (*api.Processors, error) {
	all := make(api.Processors, 0)
	for _, processor := range p.base {
		all = append(all, processor.Metadata())
	}
	return &all, nil
}

func (p ProcessorExecutor) GetProcessor(procID string) (*api.Processor, error) {
	processor, ok := p.base[procID]
	if !ok {
		return nil, fmt.Errorf("processor %s not found", procID)
	} else {
		p := processor.Metadata()
		return &p, nil
	}
}

func (p ProcessorExecutor) GetRuns(procID string) (*api.Runs, error) {
	err := checkProcessorExists(p, procID)
	if err != nil {
		return nil, err
	}
	processor, ok := p.runMap[procID]
	if !ok {
		// no runs are registered
		return &api.Runs{}, nil
	} else {
		all := make(api.Runs, len(processor))
		for _, run := range processor {
			all = append(all, run.run)
		}
		return &all, nil
	}
}

func (p ProcessorExecutor) GetRun(procID, runID string) (*api.Run, error) {
	r, err := getRun(p, procID, runID)
	if err != nil {
		return nil, err
	}

	return &r.run, nil
}

func (p ProcessorExecutor) GetConfig(procID, runID string) (*api.Configuration, error) {
	r, err := getRun(p, procID, runID)
	if err != nil {
		return nil, err
	}
	return &r.run.Configuration, nil
}

func (p ProcessorExecutor) Configure(procID, runID string, config api.Configuration) error {
	baseProcessor, ok := p.base[procID]
	if !ok {
		return fmt.Errorf("baseProcessor %s not found", procID)
	}
	if p.runMap[procID] == nil {
		//return nil, fmt.Errorf("failed")
		log.Printf("proc %s did not have any runs, creating", procID)
		p.runMap[procID] = make(map[string]*runProcessor)
	}

	err := baseProcessor.Configure(config.Processor)
	if err != nil {
		return err
	}

	//The run won't exist here at this point, so we create it:
	rp := &runProcessor{run: api.Run{
		RunID:         runID,
		Configuration: config,
		CurrentState:  api.StateCONFIGURED,
	}, baseProcessor: baseProcessor}

	p.runMap[procID][runID] = rp
	return nil
}

func (p ProcessorExecutor) GetRunData(procID, runID string) (*api.DataSpec, error) {
	r, err := getRun(p, procID, runID)
	if err != nil {
		return nil, err
	}
	ds := api.DataSpec{
		In:  make([]api.Data, 0),
		Out: make([]api.Data, 0),
	}
	r.dataLock.RLock()
	for _, v := range r.workers {
		ds.In = append(ds.In,v.in)
		ds.Out = append(ds.Out,v.out)
	}
	r.dataLock.RUnlock()
	return &ds, nil
}

func (p ProcessorExecutor) AddData(procID, runID string, data api.Data) error {
	r, err := getRun(p, procID, runID)
	if err != nil {
		return err
	}
	setRunState(r, api.StateRUNNING)
	data.State = api.DataChunkStateWAITING
	go rManager(data, r)
	return nil
}

func NewProcessorExecutor() ProcessorExecutor {
	initProc := processors.ConfiguredProcessorRegistry().Processors
	spew.Dump(initProc)
	return ProcessorExecutor{version: Version, base: initProc, runMap: make(runProcMap)}
}
