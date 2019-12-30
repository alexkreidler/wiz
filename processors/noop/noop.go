package noop

import (
	"github.com/alexkreidler/wiz/api"
	"github.com/alexkreidler/wiz/processors/processor"
	"log"
	"time"
)

type NoopProcessor struct {
	state chan api.DataChunkState
	data  interface{}
}

func (n *NoopProcessor) Configure(config interface{}) error {
	return nil
}

func (n *NoopProcessor) GetConfig() interface{} {
	return nil
}

func (n *NoopProcessor) New() processor.ChunkProcessor {
	log.Println("Creating new", n.Metadata().Name, "processor")
	return &NoopProcessor{state: make(chan api.DataChunkState)}
}

func (n *NoopProcessor) State() <-chan api.DataChunkState {
	return n.state
}

func (n *NoopProcessor) Output() interface{} {
	return n.data
}

func (n *NoopProcessor) updateState(state api.DataChunkState) {
	n.state <- state
}

func (n *NoopProcessor) done() {
	close(n.state)
}

func (n *NoopProcessor) Run(data interface{}) {
	n.updateState(api.DataChunkStateRUNNING)

	// Setting data on one potential thread, getting on another?? no-- because Output() is only called after state is succeeded. No mutex needed!
	n.data = data
	//	DO work, maybe sleep for a bit
	time.Sleep(8 * time.Second)

	n.updateState(api.DataChunkStateSUCCEEDED)
	n.done()
}

func (n *NoopProcessor) GetError() error {
	return nil
}

func (n *NoopProcessor) Metadata() api.Processor {
	return api.Processor{
		ProcID:  "noop",
		Name:    "No Operation Processor",
		Version: "0.1.0",
	}
}
