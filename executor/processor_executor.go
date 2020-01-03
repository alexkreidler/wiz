package executor

import (
	"fmt"
	"github.com/alexkreidler/deepcopy"
	"github.com/alexkreidler/wiz/api/processors"
	"github.com/alexkreidler/wiz/processors/registration"
	procApi "github.com/alexkreidler/wiz/processors/simpleprocessor"
	"github.com/mitchellh/mapstructure"
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
	run processors.Run

	// dataLock locks the state of all the data chunks/workers
	dataLock sync.RWMutex
	// workers is a map from ChunkID to worker
	workers map[string]*Worker

	// a WaitGroup doesn't make sense cause it needs the processor goroutine which we don't control to cal  Done. We could just wrap it?? in an anon goroutine

	numCompleted uint32
}

type Worker struct {
	p   procApi.ChunkProcessor
	in  processors.Data
	out processors.Data
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

// TODO: figure out how to configure logging for the processor executor

func (p ProcessorExecutor) GetAllProcessors() (*processors.Processors, error) {
	all := make(processors.Processors, 0)
	for _, processor := range p.base {
		all = append(all, processor.Metadata())
	}
	return &all, nil
}

func (p ProcessorExecutor) GetProcessor(procID string) (*processors.Processor, error) {
	processor, ok := p.base[procID]
	if !ok {
		return nil, fmt.Errorf("processor %s not found", procID)
	} else {
		p := processor.Metadata()
		return &p, nil
	}
}

func (p ProcessorExecutor) GetAllRuns(procID string) (*processors.Runs, error) {
	err := checkProcessorExists(p, procID)
	if err != nil {
		return nil, err
	}
	processor, ok := p.runMap[procID]
	if !ok {
		// no runs are registered
		return &processors.Runs{}, nil
	} else {
		all := make(processors.Runs, len(processor))
		n := 0
		for _, run := range processor {
			all[n] = run.run
			n++
		}
		return &all, nil
	}
}

func (p ProcessorExecutor) GetRun(procID, runID string) (*processors.Run, error) {
	r, err := getRun(p, procID, runID)
	if err != nil {
		return nil, err
	}

	return &r.run, nil
}

//GetConfig must be called on an existing run
func (p ProcessorExecutor) GetConfig(procID, runID string) (*processors.Configuration, error) {
	r, err := getRun(p, procID, runID)
	if err != nil {
		return nil, err
	}
	//var baseConfig interface{}
	// checks if the run has any processor configuration, else fetches default configuration from processor itself
	//if r.run.Configuration.Processor == nil {
	//	baseConfig = r.baseProcessor.GetConfig()
	//	log.Printf("Got base configuration: %#+v", baseConfig)
	//	r.run.Configuration.Processor = baseConfig
	//}
	return &r.run.Configuration, nil
}

func (p ProcessorExecutor) Configure(procID, runID string, config processors.Configuration) error {
	baseProcessor, ok := p.base[procID]
	if !ok {
		return fmt.Errorf("baseProcessor %s not found", procID)
	}
	if p.runMap[procID] == nil {
		//return nil, fmt.Errorf("failed")
		log.Printf("proc %s did not have any runs, creating", procID)
		p.runMap[procID] = make(map[string]*runProcessor)
	}

	procConfig, err := configure(baseProcessor, config.Processor)
	config.Processor = procConfig
	//err := baseProcessor.Configure(config.Processor)
	if err != nil {
		return err
	}

	//The run won't exist here at this point, so we create it:
	rp := &runProcessor{run: processors.Run{
		RunID:         runID,
		Configuration: config,
		CurrentState:  processors.StateCONFIGURED,
	}, baseProcessor: baseProcessor}

	p.runMap[procID][runID] = rp
	return nil
}

func configure(processor procApi.ChunkProcessor, userConfig interface{}) (interface{}, error) {
	baseConfig := processor.GetConfig()
	log.Printf("Got base configuration: %#+v", baseConfig)
	log.Printf("Got user config: %#+v", userConfig)

	//baseConfig is a struct that contains the default configuration

	//userConfig is a map that contains the user configuration

	// this function abstracts the merging functionality so it is the same for all projects
	// It will Configure the Processor with the same type struct options as it got from GetConfig() (aka the default config)

	//Goals: take the interface{} baseConfig which contains a struct{} options
	// and decode a map[string]interface{} into that

	// See test/main.go for a standalone example of this

	// First we need to copy the underlying struct

	bc := deepcopy.Copy(baseConfig, deepcopy.Options{ReturnPointer: true})

	err := mapstructure.Decode(userConfig, &bc)
	if err != nil {
		log.Println(err)
	}

	log.Printf("New config to apply: %#+v", bc)
	return bc, processor.Configure(bc)
}

func (p ProcessorExecutor) GetData(procID, runID string) (*processors.DataSpec, error) {
	r, err := getRun(p, procID, runID)
	if err != nil {
		return nil, err
	}
	ds := processors.DataSpec{
		In:  make([]processors.Data, 0),
		Out: make([]processors.Data, 0),
	}
	r.dataLock.RLock()
	for _, v := range r.workers {
		ds.In = append(ds.In, v.in)
		ds.Out = append(ds.Out, v.out)
	}
	r.dataLock.RUnlock()
	return &ds, nil
}

func (p ProcessorExecutor) AddData(procID, runID string, data processors.Data) error {
	r, err := getRun(p, procID, runID)
	if err != nil {
		return err
	}
	setRunState(r, processors.StateRUNNING)
	data.State = processors.DataChunkStateWAITING
	go rManager(data, r)
	return nil
}

func NewProcessorExecutor() ProcessorExecutor {
	initProc := registration.ConfiguredProcessorRegistry().Processors
	//spew.Dump(initProc)
	return ProcessorExecutor{version: Version, base: initProc, runMap: make(runProcMap)}
}
