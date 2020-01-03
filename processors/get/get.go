/*
The get processor uses go-getter from Hashicorp to quickly and easily get data from many sources including Git, HTTP, S3, etc
It takes in a simple string as input which contains a reference to the file and also any required configuration such as S3 access keys or the Git SHA
*/
package get

import (
	"github.com/alexkreidler/wiz/api/processors"
	"github.com/alexkreidler/wiz/processors/simpleprocessor"
	gogetter "github.com/hashicorp/go-getter"
	"github.com/imdario/mergo"
	"github.com/mitchellh/mapstructure"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type GoGetConfig struct {
	// Source is the source to download. It can either be a file or a folder, and Go-Getter will fetch it with GetAny
	// It can also include any go-getter configuration
	Source string
}

type GetProcessor struct {
	state  chan processors.DataChunkState
	config GoGetConfig //TODO: change this to your config type
	dir    string
}

func (p *GetProcessor) Configure(config interface{}) error {
	opts := config.(*GoGetConfig)
	p.config = *opts

	return nil
}

func (p *GetProcessor) GetConfig() interface{} {
	return p.config
}

func (p *GetProcessor) New() simpleprocessor.ChunkProcessor {
	log.Println("Creating new", p.Metadata().Name, "processor")
	return &GetProcessor{state: make(chan processors.DataChunkState)}
}

func (p *GetProcessor) State() <-chan processors.DataChunkState {
	return p.state
}

func (p *GetProcessor) Output() interface{} {
	return map[string]string{"Dir": "file://" + p.dir}
}

func (p *GetProcessor) updateState(state processors.DataChunkState) {
	p.state <- state
}

func (p *GetProcessor) done() {
	close(p.state)
}

type prog struct {
	Logger *log.Logger
}

func newProg(l *log.Logger) prog {
	p := prog{Logger:l}
	p.Logger.Println("Test")
	return p
}

func (p prog) TrackProgress(src string, currentSize, totalSize int64, stream io.ReadCloser) (body io.ReadCloser) {
	p.Logger.Println(src, currentSize, totalSize)
	return stream
}

func (p *GetProcessor) get(source string) {
	dir, err := ioutil.TempDir("", "go-get")
	if err != nil {
		log.Println(err)
		p.updateState(processors.DataChunkStateFAILED)
		return
	}
	log.Println("created temp dir", dir)

	progress := newProg(log.New(os.Stdout, "", 0))

	// TODO: expose more options
	err = gogetter.GetAny(dir, source, gogetter.WithProgress(progress))
	if err != nil {
		log.Println(err)
		p.updateState(processors.DataChunkStateFAILED)
		return
	}
	p.dir = dir
	return
}

func (p *GetProcessor) Run(data interface{}) {
	defer p.done()
	p.updateState(processors.DataChunkStateRUNNING)

	// First we decode the map into the correct structure
	var opts GoGetConfig
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

	p.get(opts.Source)

	p.updateState(processors.DataChunkStateSUCCEEDED)
	//p.done()
}

func (p *GetProcessor) GetError() error {
	return nil
}

func (p *GetProcessor) Metadata() processors.Processor {
	return processors.Processor{
		ProcID:  "get",
		Name:    "Go-Getter (Hashicorp) Processor",
		Version: "0.1.0",
	}
}
