/*
This package defines the Wiz Processor API
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

// Run is an instance of a processor associated with a task graph
type Run struct {
	RunID string
	Config Configuration // `json:"config"`
	CurState State //`json:"state"`
}
// Runs is an unordered set of runs
type Runs []Run

// Configuration is a generic type for processor-specific configuration
type Configuration interface{}

// Data is a data chunk
// In the future we may extend this to include data streams
type Data struct {
	ChunkID string
	Type DataType
	State DataChunkState

	RawData interface{}
	FilesystemReference

	OutputChunkID string
}

// DataSpec defines both the input and output data chunks in processor
type DataSpec struct {
	In []Data
	Out []Data
}

// FilesystemReference is a reference to either a file or directory
type FilesystemReference struct {
	Driver string // the filesystem driver (e.g. NFS, local, ZFS, etc)
	Location string // the actual file path location
}
