package noop

import (
	"github.com/alexkreidler/wiz/api"
	"github.com/alexkreidler/wiz/processors/processor"
	"github.com/davecgh/go-spew/spew"
	"time"
)

type NoopProcessor struct {
	state    chan api.State
	curState api.State
}

func (n NoopProcessor) New(config interface{}) (processor.Processor, error) {
	spew.Dump("got:")
	spew.Dump(config)
	return NoopProcessor{state: make(chan api.State), curState: api.State_CONFIGURED}, nil
}

func (n NoopProcessor) State() chan api.State {
	n.state <- n.curState
	return n.state
}

func (n NoopProcessor) updateState(state api.State) {
	n.state <- state
	n.curState = state
}

func (n NoopProcessor) Run(data interface{}) {
	n.updateState(api.State_RUNNING)

	//	DO work, maybe sleep for a bit
	time.Sleep(2 * time.Second)
	n.updateState(api.State_SUCCESS)
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
