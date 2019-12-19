package processors

import (
	"github.com/alexkreidler/wiz/processors/processor"
	"github.com/alexkreidler/wiz/processors/noop"
)

//
//var DefaultProcessors = map[string]Processor{
//	"noop": noop.NoopProcessor{},
//}


type ProcessorMap map[string]processor.Processor

type ProcessorRegistry struct {
	Processors ProcessorMap
}

func (p ProcessorRegistry) AddProcessor(name string, processor processor.Processor) {
	p.Processors[name] = processor;
}

func ConfiguredProcessorRegistry() ProcessorRegistry {
	p := ProcessorRegistry{make(ProcessorMap)}
	p.AddProcessor("noop", noop.NoopProcessor{})
	return p
}