package executor

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

	// StartExecutor starts the Wiz Executor binary on the given node. TODO: figure out scheduling and what nodes actually mean.
	// For now as we're only implementing the local executor this can be anything
	StartExecutor(node string) error
}

type SerializableEnv struct {
	// Name is the canonical name of the Environment, hardcoded in each implementation
	Name string
	// Description is a human readable description that may include dynamic information from the configuration:
	// e.g. local hostname or k8s namespace
	Description string

	// Configuration contains the current state of the environment's configuration
	Configuration interface{}
}