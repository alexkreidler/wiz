package noop

import (
	"github.com/alexkreidler/wiz/api/processors"
	"github.com/alexkreidler/wiz/processors/simpleprocessor"
	"log"
	"time"
)

type NoopProcessor struct {
	state chan processors.DataChunkState
	data  interface{}
}

func (n *NoopProcessor) Configure(config interface{}) error {
	return nil
}

func (n *NoopProcessor) GetConfig() interface{} {
	return nil
}

func (n *NoopProcessor) New() simpleprocessor.ChunkProcessor {
	log.Println("Creating new", n.Metadata().Name, "processor")
	return &NoopProcessor{state: make(chan processors.DataChunkState)}
}

func (n *NoopProcessor) State() <-chan processors.DataChunkState {
	return n.state
}

func (n *NoopProcessor) Output() interface{} {
	return n.data
}

func (n *NoopProcessor) updateState(state processors.DataChunkState) {
	n.state <- state
}

func (n *NoopProcessor) done() {
	close(n.state)
}

func (n *NoopProcessor) Run(data interface{}) {
	n.updateState(processors.DataChunkStateRUNNING)

	// Setting data on one potential thread, getting on another?? no-- because Output() is only called after state is succeeded. No mutex needed!
	n.data = data
	//	DO work, maybe sleep for a bit
	time.Sleep(8 * time.Second)

	n.updateState(processors.DataChunkStateSUCCEEDED)
	n.done()
}

func (n *NoopProcessor) GetError() error {
	return nil
}

func (n *NoopProcessor) Metadata() processors.Processor {
	return processors.Processor{
		ProcID:  "noop",
		Name:    "No Operation Processor",
		Version: "0.1.0",
	}
}
