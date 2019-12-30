package tasks

import (
	"fmt"
	"github.com/alexkreidler/wiz/utils/gutils"
	"github.com/davecgh/go-spew/spew"
	"github.com/segmentio/ksuid"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/traverse"
	"log"
	"reflect"
)

type TaskGraph interface {
	graph.DirectedBuilder
}

// Pipeline represents one Wiz Tasks Framework pipeline
type Pipeline struct {
	Name     string
	Graph    TaskGraph
	rootNode graph.Node
	Data     interface{}
	Spec     PipelineSpec
}

// PipelineSpec defines how a pipeline should be structured/serialized
type PipelineSpec Children

type Sequential []ProcessorNode

type Parallel Sequential //map[string]ProcessorNode

type Children struct {
	Sequential Sequential
	Parallel   Parallel
}

// ProcessorNode represents a single ETL processor in the pipeline
// TODO: deal with data merging
type ProcessorNode struct {
	id   int64
	Name string
	// ProcID is in .Processor.ID
	RunID     string
	Processor Processor
	Children  Children
}

func (p *ProcessorNode) ID() int64 {
	return p.id
}

func (p *ProcessorNode) DOTID() string {
	return p.Name
}

func NewPipeline(name string) *Pipeline {
	p := Pipeline{Name: name, Graph: simple.NewDirectedGraph()}

	// In our graph we make sure that the ID of rootNode is 0. Useful
	p.rootNode = &ProcessorNode{id: 0, Name: name + " Pipeline: Root Node"}
	p.Graph.AddNode(p.rootNode)
	return &p
}

func processorParallel(g graph.DirectedBuilder, parentNode graph.Node, p ProcessorNode) {
	// uses the node ID allocations from the regular node, but then sets to an ID in our custom node
	n := g.NewNode()
	id := n.ID()

	node := &p
	node.id = id

	g.AddNode(node)
	e := g.NewEdge(parentNode, node)
	g.SetEdge(e)

	processChildren(g, p, node)
}
func processorSequential(g graph.DirectedBuilder, previousNode graph.Node, p ProcessorNode) graph.Node {
	n := g.NewNode()
	id := n.ID()

	node := &p
	node.id = id

	g.AddNode(node)
	e := g.NewEdge(previousNode, node)
	g.SetEdge(e)

	processChildren(g, p, node)
	return node
}

func processChildren(builder graph.DirectedBuilder, c ProcessorNode, currentNode graph.Node) {
	prev := currentNode
	for _, proc := range c.Children.Sequential {
		prev = processorSequential(builder, prev, proc)
	}
	for _, proc := range c.Children.Parallel {
		processorParallel(builder, currentNode, proc)
	}
}

// UpdateFromSpec recursively builds the internal graph representation of the task graph
func (p *Pipeline) UpdateFromSpec() {
	prev := p.rootNode
	for _, proc := range p.Spec.Sequential {
		prev = processorSequential(p.Graph, prev, proc)
	}
	for _, proc := range p.Spec.Parallel {
		processorParallel(p.Graph, p.rootNode, proc)
	}
}

//func (p Pipeline) IterateChildNodes(f func(p ProcessorNode) interface{}) {
//	for _, proc := range p.Spec.Sequential {
//		f(proc)
//	}
//	for _, proc := range p.Spec.Parallel {
//		f(proc)
//	}
//}

var castError = fmt.Errorf("failed to cast on node")

// Walk does a breadth-first traversal of the pipeline's graph starting at the root node
// interrupts and return any errors that occur
func (p Pipeline) Walk(f func(p ProcessorNode) error) (err error) {
	defer func() {
		//	handles both failure to cast to processorNode and any user-function errors
		//err = recover().(error)
	}()
	trav := traverse.BreadthFirst{
		Visit: func(node graph.Node) {
			proc, ok := node.(*ProcessorNode)
			if !ok {
				panic(castError)
			}
			err := f(*proc)
			if err != nil {
				panic(err)
			}
		},
		Traverse: nil,
	}
	trav.Walk(p.Graph, p.rootNode, nil)
	return nil
}

func hasOneTypeOf(p Pipeline, t string) bool {
	for _, proc := range p.Spec.Sequential {
		if proc.Processor.Type == t {
			return true
		}
	}
	for _, proc := range p.Spec.Parallel {
		if proc.Processor.Type == t {
			return true
		}
	}
	return false
}

// CheckValidity ensures that the pipeline has 1. an input 2. input data and 3. an output and returns an error if it doesn't. Returns nil if OK
// Do we even need this?
func (p Pipeline) CheckValidity() error {
	//if !hasOneTypeOf(p, "input") {
	//	return fmt.Errorf("pipeline %s does not have an input node", p.ID)
	//}
	//if !hasOneTypeOf(p, "output") {
	//	return fmt.Errorf("pipeline %s does not have an output node", p.ID)
	//}

	if p.Data == nil {
		return fmt.Errorf("pipeline %s does not have any input data", p.Name)
	}
	return nil
}

func (p *Pipeline) AssignRunIDs() {
	x := p.Graph.Nodes()
	log.Println(reflect.ValueOf(x).Type())
	spew.Dump(p.Graph)
	gutils.IterateChildNodes(x, func(n graph.Node) {
		if n.ID() != 0 {
			node := p.Graph.Node(n.ID())

			procNode, ok := node.(*ProcessorNode)
			if !ok {
				log.Println("failed to cast")
			}

			runID := ksuid.New().String()
			log.Printf("Assigning RunID %s for processor %s (%d)", runID, procNode.Name, n.ID())
			procNode.RunID = runID
		}
	})
}
