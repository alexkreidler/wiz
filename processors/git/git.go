package git

import (
	"github.com/alexkreidler/wiz/api"
	"github.com/alexkreidler/wiz/processors/processor"
	"github.com/imdario/mergo"
	"github.com/mitchellh/mapstructure"
	git "gopkg.in/src-d/go-git.v4"
	"io/ioutil"
	"log"
)

type GitProcessor struct {
	state  chan api.DataChunkState
	config git.CloneOptions
}

func (g GitProcessor) Configure(config interface{}) error {
	var opts git.CloneOptions
	err := mapstructure.Decode(config, opts)
	if err != nil {
		return err
	}
	g.config = opts
	return nil
}

func (g GitProcessor) New() processor.ChunkProcessor {
	log.Println("Creating new", g.Metadata().Name, "processor")
	return GitProcessor{state: make(chan api.DataChunkState)}
}

func (g GitProcessor) State() <-chan api.DataChunkState {
	return g.state
}

func (g GitProcessor) Output() interface{} {
	return map[string]string{"test": "output"}
}

func (g GitProcessor) updateState(state api.DataChunkState) {
	g.state <- state

}

func (g GitProcessor) done() {
	close(g.state)
}

func (g GitProcessor) Run(data interface{}) {
	g.updateState(api.DataChunkStateRUNNING)

	var opts git.CloneOptions
	err := mapstructure.Decode(data, &opts)
	if err != nil {
		log.Println(err)
	}
	err = mergo.Merge(opts, &g.config)
	if err != nil {
		log.Println(err)
	}
	dir, err := ioutil.TempDir("", "clone-example")
	if err != nil {
		log.Println(err)
	}
	git.PlainClone(dir, false, &opts)

	g.updateState(api.DataChunkStateSUCCEEDED)
	g.done()
}

func (g GitProcessor) GetError() error {
	return nil
}

func (g GitProcessor) Metadata() api.Processor {
	return api.Processor{
		ProcID:  "git",
		Name:    "git Clone Processor",
		Version: "0.1.0",
	}
}
