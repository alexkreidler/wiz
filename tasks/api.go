package tasks

// PipelineSpec defines how a pipeline should be structured/serialized
type PipelineSpec Children

type Sequential []*ProcessorNode

type Parallel Sequential //map[string]ProcessorNode

type Children struct {
	Sequential Sequential
	Parallel   Parallel
}

// ProcessorNode represents a single ETL processor in the pipeline
// TODO: deal with data merging
type ProcessorNode struct {
	id int64
	// The public name of the node
	Name string

	// Whether the node should receive the initial data specified in the pipeline. By default, this only gets enabled for the direct children of the root node.
	GetsInitialData bool

	// The assigned runID that this processor instance gets
	RunID string

	// ProcID is in .Processor.ID
	Processor Processor
	Children  Children
}

func (p *ProcessorNode) ID() int64 {
	return p.id
}

func (p *ProcessorNode) DOTID() string {
	return p.Name
}

