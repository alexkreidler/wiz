package registration

import (
	"github.com/alexkreidler/wiz/processors/get"
	"github.com/alexkreidler/wiz/processors/git"
	"github.com/alexkreidler/wiz/processors/ls"
	"github.com/alexkreidler/wiz/processors/noop"
	"github.com/alexkreidler/wiz/processors/simpleprocessor"
)

type ProcessorMap map[string]simpleprocessor.ChunkProcessor

type ProcessorRegistry struct {
	Processors ProcessorMap
}

func (p ProcessorRegistry) AddProcessor(name string, processor simpleprocessor.ChunkProcessor) {
	p.Processors[name] = processor
}

func ConfiguredProcessorRegistry() ProcessorRegistry {
	p := ProcessorRegistry{make(ProcessorMap)}
	p.AddProcessor("noop", &noop.NoopProcessor{})
	p.AddProcessor("git", &git.GitProcessor{})
	p.AddProcessor("get", &get.GetProcessor{})
	p.AddProcessor("ls", &ls.LsProcessor{})
	return p
}
