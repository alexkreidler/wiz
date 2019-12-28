package processors

import (
	"github.com/alexkreidler/wiz/processors/git"
	"github.com/alexkreidler/wiz/processors/processor"
	"github.com/alexkreidler/wiz/processors/noop"
)

//
//var DefaultProcessors = map[string]Processor{
//	"noop": noop.NoopProcessor{},
//}


type ProcessorMap map[string]processor.ChunkProcessor

type ProcessorRegistry struct {
	Processors ProcessorMap
}

func (p ProcessorRegistry) AddProcessor(name string, processor processor.ChunkProcessor) {
	p.Processors[name] = processor
}

func ConfiguredProcessorRegistry() ProcessorRegistry {
	p := ProcessorRegistry{make(ProcessorMap)}
	p.AddProcessor("noop", noop.NoopProcessor{})
	p.AddProcessor("git", git.GitProcessor{})
	return p
}