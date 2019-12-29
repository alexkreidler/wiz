package runner

import "github.com/alexkreidler/wiz/tasks"

type Runner interface {
	//Generally, for long-running tasks and pipelines, if we use a Manager-active design, then the manager will have to explicitly control
	// all steps, and thus be running for the lifetime of the pipeline. In a local Manager environment, this won't happen. However, since
	// we'll likely be using the local executor, and only be placing on one node, it won't be a problem.

	//The environment should be set at the pipeline level and apply to all nodes, but we may have cross-env features in the future
	Run(p tasks.Pipeline) error
}
