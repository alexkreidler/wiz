package api

// ProcessorAPI is the server API for the Wiz Processor API.
type ProcessorAPI interface {
	// GetAllProcessors lists all processors on an endpoint
	GetAllProcessors() (*Processors, error)
	GetProcessor(procID string) (*Processor, error)
	// GetAllRuns lists all runs on a processor
	GetAllRuns(procID string) (*Runs, error)
	// GetRun returns an individual run
	GetRun(procID, runID string) (*Run, error)
	// GetConfig gets the current configuration of a processor Run
	GetConfig(procID, runID string) (*Configuration, error)

	// Configure accepts configuration serialized to a []byte
	// It runs synchronously and returns an error if the configuration is rejected
	// It needs the Run ID to be unique and new (e.g. no existing run)
	// It will create a new Run with the specified configuration
	Configure(procID, runID string, config Configuration) error

	// TODO: maybe provide streaming view of processor state. In go this would be a channel, gRPC a stream, IDK about HTTP

	// GetData retrieves all of the data chunks associated with a Run
	GetData(procID, runID string) (*DataSpec, error)
	// Returns nothing on success, error if empty
	AddData(procID, runID string, data Data) error
}

var ProcessorAPIEndpoints map[string]string = map[string]string{
	"GetAllProcessors": "/processors",
	"GetProcessor":     "/processors/%s",
	"GetAllRuns":       "/processors/%s/runs",
	"GetRun":           "/processors/%s/runs/%s",
	"GetConfig":        "/processors/%s/runs/%s/config",
	"Configure":        "/processors/%s/runs/%s/config",
	"GetData":          "/processors/%s/runs/%s/data",
	"AddData":          "/processors/%s/runs/%s/data",
}
