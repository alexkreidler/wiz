package processors

import "github.com/alexkreidler/wiz/processors/noop"

var DefaultProcessors = map[string]Processor{
	"noop": noop.NoopProcessor{},
}

