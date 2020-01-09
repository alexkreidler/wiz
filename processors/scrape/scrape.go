package scrape

import (
	jsonscrape "github.com/alexkreidler/jsonscrape/lib"
	"github.com/alexkreidler/wiz/api/processors"
	"github.com/alexkreidler/wiz/processors/simpleprocessor"
	"github.com/imdario/mergo"
	"github.com/mitchellh/mapstructure"
	"log"
)

type ScrapeProcessor struct {
	state   chan processors.DataChunkState
	config  jsonscrape.Config // TODO: change this to your config type
	out     interface{}
	scraper jsonscrape.Scraper
}

func (p *ScrapeProcessor) Configure(config interface{}) error {
	// TODO: cast to your config type
	opts := config.(*jsonscrape.Config)
	p.config = *opts

	return nil
}

func (p *ScrapeProcessor) GetConfig() interface{} {
	return p.config
}

func (p *ScrapeProcessor) New() simpleprocessor.ChunkProcessor {
	log.Println("Creating new", p.Metadata().Name, "processor")
	return &ScrapeProcessor{state: make(chan processors.DataChunkState)}
}

func (p *ScrapeProcessor) State() <-chan processors.DataChunkState {
	return p.state
}

func (p *ScrapeProcessor) Output() interface{} {
	return p.out
}

func (p *ScrapeProcessor) updateState(state processors.DataChunkState) {
	p.state <- state
}

func (p *ScrapeProcessor) done() {
	close(p.state)
}

func (p *ScrapeProcessor) Run(data interface{}) {
	defer p.done()
	p.updateState(processors.DataChunkStateRUNNING)

	// Remember to decode from the map[string]interface{} data to your config type
	// First we decode the map into the correct structure
	var opts jsonscrape.Config
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

	scraper, err := jsonscrape.NewScraper(opts)
	if err != nil {
		p.updateState(processors.DataChunkStateFAILED)
		return
	}

	r, err := scraper.Scrape()
	if err != nil {
		p.updateState(processors.DataChunkStateFAILED)
		return
	}

	p.out = r

	p.updateState(processors.DataChunkStateSUCCEEDED)
}

func (p *ScrapeProcessor) GetError() error {
	return nil
}

func (p *ScrapeProcessor) Metadata() processors.Processor {
	return processors.Processor{
		ProcID:  "scrape",
		Name:    "JSONScrape web scraper",
		Version: "0.1.0",
	}
}
