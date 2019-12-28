package api

// Data is a data chunk
// In the future we may extend this to include data streams
type Data struct {
	ChunkID string
	Format  DataFormat
	Type    DataType
	State   DataChunkState

	RawData interface{}
	FilesystemReference FilesystemReference

	// this is a reference to the opposing data chunk. If it is an input chunk this references the output.
	// If it is an output chunk it references the input chunk
	AssociatedChunkID string

	// TODO: look into adding provenance data to this
	// TODO: figure out how this will work with external data streams
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
