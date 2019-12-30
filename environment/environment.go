package environment

// Environment represents a deployment/execution environment where tasks can be executed
type Environment interface {
	// Name is the human name of the environment. Examples: kubernetes, ECS, local
	Name() string
	// Configure connects to the environment with optional parameters and returns an error if unsuccessful. This may require authentication, API keys, etc
	// This can be called more than once to change which specific sub-environment to connect to. E.g. which k8s namespace
	Configure(interface{}) error
	// GetCurrentConfiguration returns the current environment configuration
	GetCurrentConfiguration() interface{}
	// IsValidConfiguration returns true if the current configuration is valid and a connection to the environment can be established
	IsValidConfiguration() bool

	// Describe returns information about the environment that can be serialized
	Describe() SerializableEnv

	// StartExecutor starts the Wiz Executor binary on the given node. TODO: figure out scheduling and what nodes actually mean.
	// For now as we're only implementing the local executor this can be anything
	// Additionally, start executor should fork the process or do anything necessary so the executor continues to run
	// It can optionally return details about the started executor such as PID or Pod ID
	StartExecutor(node string) (interface{}, error)
}

//SerializableEnv is a snapshot of the environment's state at a given time.
type SerializableEnv struct {
	// Name is the canonical name of the Environment, hardcoded in each implementation
	EnvironmentID string
	// Description is a human readable description that may include dynamic information from the configuration:
	// e.g. local hostname or k8s namespace
	Description string

	//Host references a host with a valid Processor API endpoint. This can be generated from configuration or similar but must be available
	Host string

	// Configuration contains the current state of the environment's configuration
	Configuration interface{}

	// State is the state of the environment, e.g. what objects are actually applied/the real executor
	// This usually is a k8s Pod ID or a local process ID/PID
	State interface{}
}
