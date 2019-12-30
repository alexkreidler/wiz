package local

import (
	"encoding/json"
	"fmt"
	"github.com/alexkreidler/wiz/api"
	"github.com/alexkreidler/wiz/client"
	"github.com/alexkreidler/wiz/utils/gutils"
	"github.com/segmentio/ksuid"
	"gonum.org/v1/gonum/graph"

	"github.com/alexkreidler/wiz/environment"
	"github.com/alexkreidler/wiz/environment/local"

	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/alexkreidler/wiz/tasks"
	"github.com/davecgh/go-spew/spew"
)

type Manager struct {
	storageLocation string // the location that the local manager persists state to
	State
}

func NewManager(storageLocation string) *Manager {
	return &Manager{storageLocation: storageLocation, State: State{
		Pipelines:    make(map[string]tasks.Pipeline),
		Environments: make(map[string]environment.SerializableEnv),
	}}
}

// State represents the manager state. It needs to be serializable to a file
type State struct {
	Pipelines          map[string]tasks.Pipeline
	Environments       map[string]environment.SerializableEnv
	CurrentEnvironment string
}

//func fileExists(filename string) bool {
//	info, err := os.Stat(filename)
//	if os.IsNotExist(err) {
//		return false
//	}
//	return !info.IsDir()
//}

func (l *Manager) readState() error {
	f, err := ioutil.ReadFile(l.storageLocation)
	if err != nil {
		return err
	}
	return json.Unmarshal(f, &l.State)
}

func ensureDir(fileName string) {
	dirName := filepath.Dir(fileName)
	if _, serr := os.Stat(dirName); serr != nil {
		merr := os.MkdirAll(dirName, os.ModePerm)
		if merr != nil {
			panic(merr)
		}
	}
}

func (l *Manager) writeState() error {
	dat, err := json.Marshal(l.State)
	if err != nil {
		return err
	}

	ensureDir(l.storageLocation)

	return ioutil.WriteFile(l.storageLocation, dat, 0644)
}

func (l *Manager) ResetState() error {
	return os.Remove(l.storageLocation)
}

// Starts the local executor if it hasn't been started already
func (l *Manager) maybeStartLocalEnv() error {
	if _, ok := l.Environments["local"]; !ok {
		config := local.Environment{Port: 8080}
		log.Println("Starting local environment with config", config)

		e := local.NewEnvironment()
		err := e.Configure(config)
		if err != nil {
			return err
		}

		l.Environments["local"] = e.Describe()
		l.CurrentEnvironment = "local"

		e.StartExecutor("")
	}
	return nil
}

// RN it goes many Processors --> many downstream locations, etc
// TODO: think about centralizing this into a data structure that is easier to reason about

func setupProcessor(l Manager, pipeline tasks.Pipeline, node tasks.ProcessorNode) error {
	e := l.Environments[l.CurrentEnvironment]
	if e.Host == "" {
		return fmt.Errorf("failed, invalid host")
	}

	id := node.Processor.ID
	if id == "" {
		// We skip the root node
		return nil
	}

	// Setup HTTP client
	c := client.NewClient(e.Host)

	// GET /processor/id
	// Make sure its found, return error
	_, err := c.GetProcessor(id)
	if err != nil {
		return err
	}

	log.Printf("Creating run %s for processor %s (%s)", node.RunID, node.Name, id)

	//_ = api.Configuration{}

	downstreamLocs := make([]api.DownstreamDataLocation, 0)

	gutils.IterateChildNodes(pipeline.Graph.From(node.ID()), func(n graph.Node) {
		log.Println("got node", n.ID())
		procNode, ok := n.(*tasks.ProcessorNode)
		if !ok {
			log.Println("failed to cast")
		}
		log.Println("child node", procNode.Name)
		// TODO: think about which of these procNode.procesor things should be exposed vs private
		// Maybe add a function in the tasks package which returns the DownstreamDataLocation for a given procNode

		// This assumes that all RunIDs have been assigned in advance
		downstreamLocs = append(downstreamLocs, api.DownstreamDataLocation{Hostname: e.Host, ProcID: procNode.Processor.ID, RunID: procNode.RunID})
	})

	log.Println("About to configure with downstreams:", downstreamLocs)

	// POST /proc/id/run/id/config
	// Configure with Downstream True
	return c.Configure(id, node.RunID, api.Configuration{
		ExpectedData: api.ExpectedData{
			NumChunks: 1,
		},
		ExecutorConfig: api.ExecutorConfig{
			SendDownstream:      true,
			DownstreamLocations: downstreamLocs,
		},
		Processor: node.Processor.Configuration,
	})
	return nil
}

func (l *Manager) CreatePipeline(p tasks.Pipeline, environmentName string) error {
	// TODO: read state from file
	l.readState()
	defer l.writeState()

	if _, ok := l.State.Pipelines[p.Name]; ok {
		return fmt.Errorf("pipeline already exists")
	}
	localPipeline := &p
	localPipeline.UpdateFromSpec()
	err := localPipeline.CheckValidity()
	if err != nil {
		return err
	}

	spew.Dump(p.Spec)
	log.Println("Pipeline", localPipeline.Name, "is valid, creating...")
	log.Println("Assigning runIDs to processors")

	localPipeline.AssignRunIDs()
	localPipeline.UpdateInitialDataFlags()

	err = l.maybeStartLocalEnv()
	if err != nil {
		return err
	}

	// First we setup all the nodes
	err = localPipeline.Walk(func(n tasks.ProcessorNode) error {
		return setupProcessor(*l, p, n)
	})
	if err != nil {
		log.Println(err)
	}

	// And then we provide the initial data to the nodes that need it

	// In the future we may do these in two steps, but its here to avoid a situation where the processor tries to send data to a downstream that is not yet configured
	err = localPipeline.Walk(func(n tasks.ProcessorNode) error {
		return provideInitialData(*l, p, n)
	})
	if err != nil {
		log.Println(err)
	}

	return nil
}

func provideInitialData(manager Manager, p tasks.Pipeline, n tasks.ProcessorNode) error {
	e := manager.Environments[manager.CurrentEnvironment]
	if e.Host == "" {
		return fmt.Errorf("failed, invalid host")
	}

	id := n.Processor.ID
	if id == "" {
		// We skip the root node
		return nil
	}

	if n.GetsInitialData {
		// Setup HTTP client
		c := client.NewClient(e.Host)

		chunkID := ksuid.New().String()

		log.Printf("Providing initial data to node %s (%s) with Chunk ID: %s", n.Name, n.Processor.ID, chunkID)

		return c.AddData(n.Processor.ID, n.RunID, api.Data{
			ChunkID:             chunkID,
			Format:              api.DataFormatRAW,
			Type:                api.DataTypeINPUT,
			State:               api.DataChunkStateWAITING,
			RawData:             p.Data,
			FilesystemReference: api.FilesystemReference{},
			AssociatedChunkID:   ksuid.New().String(),
		})
	}
	return nil
}

func (l Manager) ReadPipeline(name string) (tasks.Pipeline, error) {
	p, ok := l.State.Pipelines[name]
	if !ok {
		return tasks.Pipeline{}, fmt.Errorf("failed to retrieve pipeline %s", name)
	}
	return p, nil
}

func (l Manager) DeletePipeline(name string) (tasks.Pipeline, error) {
	p, ok := l.State.Pipelines[name]
	if !ok {
		return tasks.Pipeline{}, fmt.Errorf("pipeline doesn't exist")
	}
	delete(l.State.Pipelines, name)
	return p, nil
}
