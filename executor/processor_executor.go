package executor

import (
	"fmt"
	"github.com/alexkreidler/wiz/api"
	"github.com/alexkreidler/wiz/processors"
	procApi "github.com/alexkreidler/wiz/processors/processor"
	"github.com/davecgh/go-spew/spew"
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

type runProcMap map[string]map[string]*runProcessor

//ProcessorExecutor implements the ProcessorAPIServer for builtin Golang Processors. It uses channels, maps, and concurrency to parallelize by chunks
//It is designed so all operations can use a value receiver. The version an base processors don't need to be modified, but the specific runs do, which is why it is a map to pointers
//can maps be modified with value receivers? yes they can, because maps, slices, and channels are inherently mutable
//think about concurrency issues, e.g. will multiple requests result in consistent state -- map need a stronger concurrency primitive than a simple Map, e.g. locks
type ProcessorExecutor struct {
	version string
	// base maps the ID of the processor to the processor. These are all base, non-configured processors that are registered at startup
	//Their Metadata functions are the source of all processor information
	base map[string]procApi.Processor
	//runMap is a map of processor IDs to runIDs and processors
	runMap runProcMap
}

func (p ProcessorExecutor) GetAllProcessors() (*api.Processors, error) {
	all := make(api.Processors,0)
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
			all = append(all, run.r)
		}
		return &all, nil
	}
}

func (p ProcessorExecutor) GetRun(procID, runID string) (*api.Run, error) {
	r, err := getRun(p, procID, runID)
	if err != nil {
		return nil, err
	}
	return &r.r, nil
}

func (p ProcessorExecutor) GetConfig(procID, runID string) (*api.Configuration, error) {
	r, err := getRun(p, procID, runID)
	if err != nil {
		return nil, err
	}
	return &r.r.Configuration, nil
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

	//The run won't exist here at this point, so we create it:
	proc, err := baseProcessor.New(config)
	if err != nil {
		return err
	}
	fmt.Println("got here")
	rp := &runProcessor{p: proc, r: api.Run{
		RunID: runID,
		Configuration: config,
		State:  api.StateCONFIGURED,
	}}

	spew.Dump(rp)

	p.runMap[procID][runID] = rp
	return nil
}

func (p ProcessorExecutor) GetRunData(procID, runID string) (*api.DataSpec, error) {
	r, err := getRun(p, procID, runID)
	if err != nil {
		return nil, err
	}
	return &r.ds, nil
}

func (p ProcessorExecutor) AddData(procID, runID string, data api.Data) error {
	r, err := getRun(p, procID, runID)
	if err != nil {
		return err
	}
	data.State = api.DataChunkStateRUNNING
	//todo: make all these chunks concurrent
	r.ds.In = append(r.ds.In, data)

	switch data.Type {
	case api.DataTypeRAW:
		r.p.Run(data.RawData)
		break
	case api.DataTypeFILESYSTEMREF:
		r.p.Run(data.FilesystemReference)
	}
	return nil
}


func NewProcessorExecutor() ProcessorExecutor {
	initProc := processors.ConfiguredProcessorRegistry().Processors
	spew.Dump(initProc)
	return ProcessorExecutor{version: Version, base: initProc, runMap: make(runProcMap)}
}