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
    config   ConfigType //TODO: change this to your config type
}

func (p *<%= Name %>Processor) Configure(config interface{}) error {
    // TODO: cast to your config type
	opts := config.(*ConfigType)
	p.config = *opts

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
	defer p.done()
	p.updateState(api.DataChunkStateRUNNING)

    // Remember to decode from the map[string]interface{} data to your config type
    // First we decode the map into the correct structure
    var opts ConfigType
    log.Printf("got raw data: %#+v \n", data)
    err := mapstructure.Decode(data, &opts)
    if err != nil {
        log.Println(err)
    }

    // Then we merge the config into the data
    log.Printf("existing config: %#+v, new data: %#+v \n", p.config, opts)
    err = mergo.Merge(&opts, p.config, func(config *mergo.Config) {
        config.Overwrite = true
    })
    if err != nil {
        log.Println(err)
    }
    // opts now contains the merged options

	// TODO: Add your code here

	if err != nil {
	    p.updateState(api.DataChunkStateFAILED)
	    return
	}

	p.updateState(api.DataChunkStateSUCCEEDED)
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
