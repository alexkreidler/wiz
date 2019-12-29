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

func (g *GitProcessor) Configure(config interface{}) error {
	log.Printf("got config: %#+v \n", config)

	var opts git.CloneOptions
	err := mapstructure.Decode(config, &opts)
	if err != nil {
		return err
	}
	log.Printf("opts: %#+v \n", opts)
	g.config = opts

	return nil
}
func (g *GitProcessor) GetConfig() interface{} {
	return g.config
}

func (g *GitProcessor) New() processor.ChunkProcessor {
	log.Println("Creating new", g.Metadata().Name, "processor")
	return &GitProcessor{state: make(chan api.DataChunkState)}
}

func (g *GitProcessor) State() <-chan api.DataChunkState {
	return g.state
}

func (g *GitProcessor) Output() interface{} {
	return map[string]string{"test": "output"}
}

func (g *GitProcessor) updateState(state api.DataChunkState) {
	g.state <- state

}

func (g *GitProcessor) done() {
	close(g.state)
}

func (g *GitProcessor) Run(data interface{}) {
	g.updateState(api.DataChunkStateRUNNING)

	var opts git.CloneOptions
	log.Printf("got raw data: %#+v \n", data)
	err := mapstructure.Decode(data, &opts)
	if err != nil {
		log.Println(err)
	}
	log.Printf("existing config: %#+v, new data: %#+v \n", g.config, opts)
	err = mergo.Merge(&opts, g.config)
	if err != nil {
		log.Println(err)
	}
	dir, err := ioutil.TempDir("", "clone-example")
	if err != nil {
		log.Println(err)
	}
	log.Println("created temp dir", dir)
	log.Printf("Cloning with options %#+v \n", opts)
	_, err = git.PlainClone(dir, false, &opts)
	if err != nil {
		log.Println(err)
	}
	log.Println("cloned repository")

	g.updateState(api.DataChunkStateSUCCEEDED)
	g.done()
}

func (g *GitProcessor) GetError() error {
	return nil
}

func (g *GitProcessor) Metadata() api.Processor {
	return api.Processor{
		ProcID:  "git",
		Name:    "git Clone Processor",
		Version: "0.1.0",
	}
}
