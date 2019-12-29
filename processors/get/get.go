
package get

import (
	"github.com/alexkreidler/wiz/api"
	"github.com/alexkreidler/wiz/processors/processor"
	"log"
	"time"
)

type GetProcessor struct {
	state    chan api.DataChunkState
    config   interface{} //TODO: change this to your config type
}

func (p *GetProcessor) Configure(config interface{}) error {
	p.config = config
    return nil
}

func (p *GetProcessor) GetConfig() interface{} {
	return p.config
}

func (p *GetProcessor) New() processor.ChunkProcessor {
	log.Println("Creating new", p.Metadata().Name, "processor")
	return &GetProcessor{state: make(chan api.DataChunkState)}
}

func (p *GetProcessor) State() <-chan api.DataChunkState {
	return p.state
}

func (p *GetProcessor) Output() interface{} {
	return map[string]string{"test": "output"}
}

func (p *GetProcessor) updateState(state api.DataChunkState) {
	p.state <- state
}

func (p *GetProcessor) done() {
	close(p.state)
}

func (p *GetProcessor) Run(data interface{}) {
	p.updateState(api.DataChunkStateRUNNING)

	// TODO: Add your code here

	p.updateState(api.DataChunkStateSUCCEEDED)
	p.done()
}

func (p *GetProcessor) GetError() error {
	return nil
}

func (p *GetProcessor) Metadata() api.Processor {
	return api.Processor{
		ProcID:  "get",
		Name:    "Go-Get (Hashicorp) Processor",
		Version: "0.1.0",
	}
}
