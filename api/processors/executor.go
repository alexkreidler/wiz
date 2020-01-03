package processors

type DownstreamDataLocation struct {
	Hostname string
	ProcID   string
	RunID    string
}

// ExecutorConfig is the executor configuration
type ExecutorConfig struct {
	MaxWorkers     int
	SendDownstream bool
	// TODO: make API URL templating agnostic, refactor into server package
	// e.g. Server.GetProc(procID, runID) will return /processors/procID/runs/runID
	// this could be used both for generating endpoints and here, for passing outputs downstream
	DownstreamLocations []DownstreamDataLocation
}
