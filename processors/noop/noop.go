package noop

import (
	"fmt"
	"github.com/alexkreidler/wiz/api"
	"github.com/alexkreidler/wiz/processors/processor"
	"log"
	"time"
)

type NoopProcessor struct {
	state    chan api.State
	curState api.State
}

func (n NoopProcessor) New(config interface{}) (processor.Processor, error) {
	log.Println("Creating new", n.Metadata().Name, "processor with config", config)
	return NoopProcessor{state: make(chan api.State, 3), curState: api.StateCONFIGURED}, nil
}

func (n NoopProcessor) State() chan api.State {
	//n.state <- n.curState
	return n.state
}

func (n NoopProcessor) updateState(state api.State) {
	//n.state <- state

	select {
	case n.state <- state:
		fmt.Println("sent message", state)
	default:
		fmt.Println("no message sent")
	}

	n.curState = state

}

func (n NoopProcessor) Run(data interface{}) {
	n.updateState(api.StateRUNNING)

	//	DO work, maybe sleep for a bit
	time.Sleep(200 * time.Millisecond)
	n.updateState(api.StateSUCCESS)
}

func (n NoopProcessor) GetError() error {
	return nil
}

func (n NoopProcessor) Metadata() api.Processor {
	return api.Processor{
		ProcID:  "noop",
		Name:    "No Operation Processor",
		Version: "0.1.0",
	}
}
