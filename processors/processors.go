package processors

import (
	"github.com/alexkreidler/wiz/processors/get"
	"github.com/alexkreidler/wiz/processors/git"
	"github.com/alexkreidler/wiz/processors/noop"
	"github.com/alexkreidler/wiz/processors/processor"
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
	p.AddProcessor("git", &git.GitProcessor{})
	p.AddProcessor("get", &get.GetProcessor{})
	return p
}
