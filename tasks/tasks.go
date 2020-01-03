package tasks

import (
	"fmt"
	"github.com/alexkreidler/wiz/utils"
	"github.com/alexkreidler/wiz/utils/gutils"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/traverse"
	"log"
)

type TaskGraph interface {
	graph.DirectedBuilder
}

// Pipeline represents one Wiz Tasks Framework pipeline
type Pipeline struct {
	Name     string
	Graph    TaskGraph `json:"-"`
	rootNode graph.Node
	Data     interface{}
	Spec     PipelineSpec
}

func NewPipeline(name string) *Pipeline {
	p := Pipeline{Name: name, Graph: simple.NewDirectedGraph()}

	p.addRootNode(name + " root node")
	return &p
}


func (p *Pipeline) setupRequired() {
	if p.Graph == nil {
		p.Graph = simple.NewDirectedGraph()
	}
	if p.Name == "" {
		p.Name = "unnamed"
	}
	p.addRootNode(p.Name + " root node")
}

func (p *Pipeline) addRootNode(name string) {
	if p.rootNode == nil {
		// In our graph we make sure that the ID of rootNode is 0. Useful
		p.rootNode = &ProcessorNode{id: 0, Name: name}
		p.Graph.AddNode(p.rootNode)
	}
}

func processorParallel(g graph.DirectedBuilder, parentNode graph.Node, p *ProcessorNode) {
	// uses the node ID allocations from the regular node, but then sets to an ID in our custom node
	n := g.NewNode()
	id := n.ID()

	p.id = id

	g.AddNode(p)
	e := g.NewEdge(parentNode, p)
	g.SetEdge(e)

	processChildren(g, p, p)
}
func processorSequential(g graph.DirectedBuilder, previousNode graph.Node, p *ProcessorNode) graph.Node {
	n := g.NewNode()
	id := n.ID()

	p.id = id

	g.AddNode(p)
	e := g.NewEdge(previousNode, p)
	g.SetEdge(e)

	processChildren(g, p, p)
	return p
}

func processChildren(builder graph.DirectedBuilder, c *ProcessorNode, currentNode graph.Node) {
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
	p.setupRequired()
	prev := p.rootNode
	for _, proc := range p.Spec.Sequential {
		prev = processorSequential(p.Graph, prev, proc)
	}
	for _, proc := range p.Spec.Parallel {
		processorParallel(p.Graph, p.rootNode, proc)
	}
}

var castError = fmt.Errorf("failed to cast on node")

// Walk does a breadth-first traversal of the pipeline's graph starting at the root node
// interrupts and return any errors that occur
// TODO: could just replace with IterateChildren(p.Nodes). I guess BFS can have some specific use-cases, not really for us though
func (p Pipeline) Walk(f func(p ProcessorNode) error) (err error) {
	defer func() {
		//	handles both failure to cast to processorNode and any user-function errors
		//err = recover().(error)
		//	TODO: is this idiomatic?
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

func (p *Pipeline) AssignRunIDs(overwrite bool) {
	gutils.IterateChildNodes(p.Graph.Nodes(), func(n graph.Node) {
		if n.ID() != 0 {
			procNode, ok := n.(*ProcessorNode)
			if !ok {
				log.Println("failed to cast")
			}

			runID := utils.GenID()
			log.Printf("Assigning RunID %s for processor %s (%d)", runID, procNode.Name, n.ID())

			// don't overwrite if it already exists
			if procNode.RunID != "" && !overwrite {
				return
			}
			procNode.RunID = runID
		}
	})
}

func (p *Pipeline) UpdateInitialDataFlags() {
	gutils.IterateChildNodes(p.Graph.From(0), func(n graph.Node) {
		procNode, ok := n.(*ProcessorNode)
		if !ok {
			log.Println("failed to cast")
		}
		log.Printf("Setting GetsInitialData true for processor %s (%d)", procNode.Name, n.ID())
		procNode.GetsInitialData = true
	})
}
