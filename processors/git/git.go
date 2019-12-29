package git

import (
	"github.com/alexkreidler/mergo"
	"github.com/alexkreidler/wiz/api"
	"github.com/alexkreidler/wiz/processors/processor"
	"github.com/mitchellh/mapstructure"
	git "gopkg.in/src-d/go-git.v4"
	"io/ioutil"
	"log"
)

type GitProcessor struct {
	state  chan api.DataChunkState
	config git.CloneOptions
	dir    string
}

func (g *GitProcessor) Configure(config interface{}) error {
	opts := config.(*git.CloneOptions)
	g.config = *opts

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
	return map[string]string{"Dir": g.dir}
}

func (g *GitProcessor) updateState(state api.DataChunkState) {
	g.state <- state

}

func (g *GitProcessor) done() {
	close(g.state)
}

func (g *GitProcessor) Run(data interface{}) {
	g.updateState(api.DataChunkStateRUNNING)

	// First we decode the map into the correct structure
	var opts git.CloneOptions
	log.Printf("got raw data: %#+v \n", data)
	err := mapstructure.Decode(data, &opts)
	if err != nil {
		log.Println(err)
	}

	// Then we merge the config into the data
	log.Printf("existing config: %#+v, new data: %#+v \n", g.config, opts)
	err = mergo.Merge(&opts, g.config, func(config *mergo.Config) {
		config.Overwrite = true
	})
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
	g.dir = dir

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
