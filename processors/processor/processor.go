// This package provides the Processor interface which is implemented by all the other packages in the parent package
package processor

import (
	"github.com/alexkreidler/wiz/api"
)

//Processor is a simple system for Go builtin processors. It should abstract away Runs, etc and only need configuration and data
//Each chunk should be able to spawn a new processor. In the future, this may change
//Also, each processor should only be configured once, for that chunk, thus any reconfiguration logic is not implemented yet
type Processor interface {
	// Configure configures the processor with the given configuration or an error if its invalid
	// This should not depend on any previously set configuration or data
	Configure(config interface{}) error

	// GetConfig returns the current configuration
	GetConfig() interface{}

	// Configure allows for reconfiguration of the processor, but it should provide a New method that does it automatically?
	//Configure(config interface{}) error

	Metadata() api.Processor

	//The processor should automatically shutdown any resources used by it upon the
	//sending of the Succeeded or Failed states? it will be GCd automatically
}

// ChunkProcessor is a processor with the additional restriction that it can only processes one chunk
// We may add a function to reset the processor state and process new chunks, but this would only be
// necessary if setting up the processor was computationally expensive. (e.g. memory copies etc)
// For now, it must be disposed of on completion

// For now, most builtin processors will only emit these states: WAITING, RUNNING, and succeeded/failed
type ChunkProcessor interface {
	Processor

	// New returns a new copy of the chunkProcessor based on its current configuration
	New() ChunkProcessor

	//It will be in Syncing by default, and then for our simple implementation case, go into validating, running, and then success or failure
	// It returns a receive only channel, as we can't send on to the state
	// In the future, DataChunkState may be a struct that contains metadata like progress (e.g. download progress). This is a #1 reason to use channels, for allowing rich state updates
	// If the state is only a sequence from Waiting to Validating to Determining to Running to a terminated state, then we can simply have functions which represent these things
	// and run them synchronously, having a manager update the higher up state
	State() <-chan api.DataChunkState

	//we may simply use functions
	//Validate(data)

	// This provides the raw data to the processor, either in FS form or raw form.
	// It should run the entirety of the processor's task/operation synchronously
	// Expect this function to be called as a new goroutine
	// It should update the state appropriately
	Run(data interface{})

	// Output fetches the data output on successful completion
	Output() interface{}

	//GetError returns the error which occurred if it goes from running to failed
	GetError() error
}