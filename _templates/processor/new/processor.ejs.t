---
to: processors/<%= name %>/<%= name %>.go
---

package <%= name %>

import (
	"github.com/alexkreidler/wiz/api"
	"github.com/alexkreidler/wiz/processors/processor"
	"log"
	"time"
)

type <%= Name %>Processor struct {
	state    chan api.DataChunkState
    config   interface{} //TODO: change this to your config type
}

func (p *<%= Name %>Processor) Configure(config interface{}) error {
	p.config = config
    return nil
}

func (p *<%= Name %>Processor) GetConfig() interface{} {
	return p.config
}

func (p *<%= Name %>Processor) New() processor.ChunkProcessor {
	log.Println("Creating new", p.Metadata().Name, "processor")
	return &<%= Name %>Processor{state: make(chan api.DataChunkState)}
}

func (p *<%= Name %>Processor) State() <-chan api.DataChunkState {
	return p.state
}

func (p *<%= Name %>Processor) Output() interface{} {
	return map[string]string{"test": "output"}
}

func (p *<%= Name %>Processor) updateState(state api.DataChunkState) {
	p.state <- state
}

func (p *<%= Name %>Processor) done() {
	close(p.state)
}

func (p *<%= Name %>Processor) Run(data interface{}) {
	p.updateState(api.DataChunkStateRUNNING)

	// TODO: Add your code here

	p.updateState(api.DataChunkStateSUCCEEDED)
	p.done()
}

func (p *<%= Name %>Processor) GetError() error {
	return nil
}

func (p *<%= Name %>Processor) Metadata() api.Processor {
	return api.Processor{
		ProcID:  "<%= name %>",
		Name:    "<%= description %>",
		Version: "<%= version %>",
	}
}
