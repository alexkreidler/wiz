package tasks

// Processor represents a Wiz Tasks node which can process data. This is how it is serialized for input
// It will actually be mapped to the processor specified by ID and Version and a generated RunID
// TODO: maybe make separate structs for the internal Tasks framework representation and the ones that are serialized from YAML
//e.g. ID in the lib and Name in Yaml
type Processor struct {
	ID            string // the unique name for the processor
	Version       string // the semantic version of the processor if required
	Type          string // the category of the processor: either input, output, or transformation - nil means transform
	Configuration interface{}
	// we may add more information like runtime, etc and other metadata
	// where to store the actual executor-specific information like which pod is running on?
}
