// Ls is a processor to return a list of FileInfo for the provided directory
package ls

import (
	"github.com/alexkreidler/wiz/api"
	"github.com/alexkreidler/wiz/processors/processor"
	"github.com/mitchellh/mapstructure"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type FileInfo struct {
	Name    string
	Size    int64
	Mode    os.FileMode
	ModTime time.Time
	IsDir   bool
}

type LsProcessor struct {
	state chan api.DataChunkState
	files []FileInfo
}

func (p *LsProcessor) Configure(config interface{}) error {
	return nil
}

func (p *LsProcessor) GetConfig() interface{} {
	return nil
}

func (p *LsProcessor) New() processor.ChunkProcessor {
	log.Println("Creating new", p.Metadata().Name, "processor")
	return &LsProcessor{state: make(chan api.DataChunkState)}
}

func (p *LsProcessor) State() <-chan api.DataChunkState {
	return p.state
}

func (p *LsProcessor) Output() interface{} {
	return p.files
}

func (p *LsProcessor) updateState(state api.DataChunkState) {
	p.state <- state
}

func (p *LsProcessor) done() {
	close(p.state)
}

func (p *LsProcessor) Run(data interface{}) {
	defer p.done()
	p.updateState(api.DataChunkStateRUNNING)

	// Remember to decode from the map[string]interface{} data to your config type
	// First we decode the map into the correct structure
	var opts struct {
		Dir string
	}
	log.Printf("got raw data: %#+v \n", data)
	err := mapstructure.Decode(data, &opts)
	if err != nil {
		log.Println(err)
	}

	// opts now contains the options

	files, err := ioutil.ReadDir(opts.Dir)
	if err != nil {
		log.Fatal(err)
	}

	fs := make([]FileInfo, 0)

	for _, file := range files {
		fs = append(fs, FileInfo{
			Name:    file.Name(),
			Size:    file.Size(),
			Mode:    file.Mode(),
			ModTime: file.ModTime(),
			IsDir:   file.IsDir(),
		})
	}
	p.files = fs
	// TODO: Add your code here

	if err != nil {
		p.updateState(api.DataChunkStateFAILED)
		return
	}

	p.updateState(api.DataChunkStateSUCCEEDED)
}

func (p *LsProcessor) GetError() error {
	return nil
}

func (p *LsProcessor) Metadata() api.Processor {
	return api.Processor{
		ProcID:  "ls",
		Name:    "List File Processor",
		Version: "0.1.0",
	}
}
