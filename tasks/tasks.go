package tasks

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

// Pipeline represents one Wiz Tasks Framework pipeline
type Pipeline struct {
	Name string
	g graph.DirectedBuilder
	rootNode graph.Node
	Spec PipelineSpec
}

// PipelineSpec defines how a pipeline should be structured/serialized
type PipelineSpec Children

type Sequential []ProcessorNode

type Parallel Sequential //map[string]ProcessorNode

type Children struct {
	Sequential
	Parallel
}

type ProcessorNode struct {
	id int64
	Name string
	Processor
	Children
}

func (p ProcessorNode) ID() int64 {
	return p.id
}

func (p ProcessorNode) DOTID() string {
	return p.Name
}

func NewPipeline(name string) *Pipeline {
	p := Pipeline{Name: name, g: simple.NewDirectedGraph()}
	p.rootNode = ProcessorNode{id: 0, Name: name + " Pipeline: Root Node"}
	p.g.AddNode(p.rootNode)
	return &p
}

func processorParallel(g graph.DirectedBuilder, parentNode graph.Node, p ProcessorNode) {
	n := g.NewNode()
	id := n.ID()

	node := p
	node.id = id

	g.AddNode(node)
	e := g.NewEdge(parentNode, node)
	g.SetEdge(e)

	processChildren(g, p, node)
}
func processorSequential(g graph.DirectedBuilder, previousNode graph.Node, p ProcessorNode) graph.Node {
	n := g.NewNode()
	id := n.ID()

	node := p
	node.id = id

	g.AddNode(node)
	e := g.NewEdge(previousNode, node)
	g.SetEdge(e)

	processChildren(g, p, node)
	return node
}

func processChildren(builder graph.DirectedBuilder, c ProcessorNode, currentNode graph.Node) {
	prev := currentNode
	for _, proc := range c.Sequential {
		prev = processorSequential(builder, prev, proc)
	}
	for _, proc := range c.Parallel {
		processorParallel(builder, currentNode, proc)
	}
}

// UpdateFromSpec recursively builds the internal graph representation of the task graph
func (p *Pipeline) UpdateFromSpec() {
	prev := p.rootNode
	for _, proc := range p.Spec.Sequential {
		prev = processorSequential(p.g, prev, proc)
	}
	for _, proc := range p.Spec.Parallel {
		processorParallel(p.g, p.rootNode, proc)
	}
}