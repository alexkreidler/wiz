package processors

import (
	"github.com/alexkreidler/wiz/api"
)

//Processor is a simple system for Go builtin processors. It should abstract away Runs, etc and only need configuration and data
//Each chunk should be able to spawn a new processor. In the future, this may change
//Also, each processor should only be configured once, for that chunk, thus any reconfiguration logic is not implemented yet
type Processor interface {
	// New returns a new processor with the given configuration or an error if its invalid
	New(config interface{}) (Processor, error)

	// Configure allows for reconfiguration of the processor, but it should provide a New method that does it automatically?
	//Configure(config interface{}) error

	//It will be in Syncing by default, and then for our simple implementation case, go into validating, running, and then success or failure
	State() chan api.State

	// This provides the raw data to the processor, either in FS form or raw form.
	Run(data interface{})

	//GetError returns the error which occured if it goes from validating to configured, running to configured or failed
	GetError() error

	Metadata() api.Processor

	//The processor should automatically shutdown any resources used by it upon the
	//sending of the Succeeded or Failed states
}
