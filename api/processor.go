/*
This package defines the Wiz Processor API
TODO: maybe think about moving this to a separate package as we are going to add many different API types
*/
package api

// Processor is a logical processor
type Processor struct {
	// ProcID is the uniquely identifiable Processor ID
	ProcID string
	// Name is the human readable name
	Name    string
	Version string
}

// Processors is an unordered set of processors
type Processors []Processor

// ExpectedData is the data expected by a Run
// for now it is just the number of chunks, but in the future could contain a list of ChunkIDs
// If ChunkIDs were hashes this could be a defacto form of externalized/internal state?
type ExpectedData struct {
	NumChunks uint32 // a counter value
}

// Run is an instance of a processor associated with a task graph
type Run struct {
	RunID         string
	Configuration Configuration

	// Note: Embedding structs will automatically promote the child struct's functions,
	// and since our State type is an enum that overrides the default Marshal and Unmarshal functions,
	// it overwrites it for the parent type as well.
	// Remember, the CurrentState must be updated from the RunProcessor state to be fresh. TODO: think about these guarantees
	CurrentState State
}

// Runs is an unordered set of runs
type Runs []Run

// Configuration is a generic type for processor-specific configuration
type Configuration struct {
	// Embeded structs have a weird JSON serialization issue
	ExpectedData   ExpectedData
	ExecutorConfig ExecutorConfig
	Processor      interface{}
}
