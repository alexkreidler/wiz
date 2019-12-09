/*
Package manager represents the Wiz Manager library. It exposes a Go API, and is designed to work either:

1. running locally as part of the Wiz CLI, persisting to a local storage medium

2. running as a long-running daemon that serves API requests from various Wiz CLIs (this needs auth in the future) 

For now, we are only implementing method 1
*/
package manager

import (
	"github.com/alexkreidler/wiz/executor"
	"github.com/alexkreidler/wiz/tasks"
)

// Manager abstracts the Wiz Manager library/daemon away from the CLI. When in single-user mode, the Manager interface runs in the Wiz binary and persists its own state to disk. When in multi-user mode, each function calls an HTTP API to the long-running Manager Daemon
type Manager interface {
	// CreatePipeline creates and runs a pipeline on the given executor. It returns an error if the pipeline name is not unique, or the given executor is not registered
	// In the future, we may add features that allow running pipelines on multiple executor environments. Additionally, we may separate the existence of a pipeline on the manager
	// to the actual execution of that pipeline to enable advanced features like multi-user collaboration on pipelines or multiple runs of the same pipeline with different data sources.
	CreatePipeline(p tasks.Pipeline, environment executor.Environment) error

	// ReadPipeline reads the pipeline information given a pipeline name, and returns an error if not found
	ReadPipeline(name string) (tasks.Pipeline, error)

	// DeletePipeline deletes the pipeline with a given name, returning an error if not found
	DeletePipeline(name string) tasks.Pipeline

	// CreateEnvironment adds the given already configured, VALID environment to the Manager. IsValidEnvironment should be called immediately and an error returned if it is not valid
	// In the future, for multi-user mode, in the use case that the Wiz Manager server has a service account e.g with the proper permissions and that all user
	// connections should use that configuration, then we may expose the configuration of the Environment through the Manager so that the user can choose from the CLI the env configs to use from the Manager's perspective
	CreateEnvironment(environment executor.Environment) error
	//	How environments (an interface are serialized across the CLI-manager connection (An HTTP API) will be complicated
	//	Some custom environments will only be registered with the manager?
	// ReadEnvironment returns the environment configuration for a given env name
	ReadEnvironment(name string) interface{}
	//UpdateEnvironment(name string)
	// DeleteEnvironment deletes the environment with a given name and returns an error if it doesn't exist
	DeleteEnvironment(name string) error
}