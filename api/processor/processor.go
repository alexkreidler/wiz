/*
This package defines the Wiz Processor API
 */
package processor

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

// Run is an instance of a processor associated with a task graph
type Run struct {
	RunID string
	Configuration
	State
}
// Runs is an unordered set of runs
type Runs []Run

// Configuration is a generic type for processor-specific configuration
type Configuration interface{}

// Data is a data chunk
// In the future we may extend this to include data streams
type Data struct {
	ChunkID string
	DataType

	RawData interface{}
	FilesystemReference

	OutputChunkID string
}

// DataSpec defines both the input and output data chunks in processor
type DataSpec struct {
	In Data
	Out Data
}

// FilesystemReference is a reference to either a file or directory
type FilesystemReference struct {
	Driver string // the filesystem driver (e.g. NFS, local, ZFS, etc)
	Location string // the actual file path location
}

// APIServer is the server API for the Wiz Processor API.
type APIServer interface {
	// GetAllProcessors lists all processors on an endpoint
	GetAllProcessors() (Processors, error)
	GetProcessor(procID string) (Processor, error)
	// GetRuns lists all runs on a processor
	GetRuns(procID string) (Runs, error)
	// GetRun returns an individual run
	GetRun(procID, runID string) (Run, error)
	// GetConfig gets the current configuration of a processor Run
	GetConfig(procID, runID string) (Configuration, error)

	// Configure accepts configuration serialized to a []byte
	// It runs synchronously and returns an error if the configuration is rejected
	// It needs the Run ID to be unique and new (e.g. no existing run)
	// It will create a new Run with the specified configuration
	Configure(procID, runID string, config Configuration) (error)
	
	// TODO: maybe provide streaming view of processor state. In go this would be a channel, gRPC a stream, IDK about HTTP

	// GetRunData retrieves all of the data chunks associated with a Run
	GetRunData(procID, runID string) (DataSpec, error)
	// Returns nothing on success, error if empty
	AddData(procID, runID string, data Data) (error)
}
