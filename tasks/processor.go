package tasks

// Processor represents a Wiz Tasks node which can process data. This is how it is serialized for input
// It will actually be mapped to the processor specified by Name and Version and a generated RunID
type Processor struct {
	Name          string // the unique name for the processor
	Type          string // the category of the processor: either input, output, or transformation - nil means transform
	Version       string // the semantic version of the processor if required
	Configuration interface{}
	// we may add more information like runtime, etc and other metadata
	// where to store the actual executor-specific information like which pod is running on?
}
