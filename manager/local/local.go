package local

import (
	"fmt"
	"github.com/alexkreidler/wiz/api"
	"github.com/alexkreidler/wiz/client"
	"github.com/alexkreidler/wiz/environment"
	"github.com/alexkreidler/wiz/environment/local"
	"github.com/alexkreidler/wiz/utils"
	"github.com/alexkreidler/wiz/utils/gutils"
	"github.com/davecgh/go-spew/spew"

	jsoniter "github.com/json-iterator/go"
	"gonum.org/v1/gonum/graph"
	"strings"
	"time"

	"github.com/shirou/gopsutil/process"

	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/alexkreidler/wiz/tasks"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Options struct {
	// RestartExecutor restarts the local Wiz Executor even if it is already running. Useful for development
	RestartExecutor bool
	
	// the location that the local manager persists state to
	StorageLocation string

	// DEBUG options:

	// PreserveRunIDs allows existing RunIDs in the spec to be used
	PreserveRunIDs bool
	// OverwritePipelines allows the manager to overwrite an existing pipleline
	OverwritePipelines bool
}

type Manager struct {
	Options Options
	State   State
}

func NewManager(opts Options) *Manager {
	return &Manager{State: State{
		Pipelines:    make(map[string]tasks.Pipeline),
		Environments: make(map[string]environment.SerializableEnv),
	}, Options: opts}
}

// State represents the manager state. It needs to be serializable to a file
type State struct {
	Pipelines          map[string]tasks.Pipeline
	Environments       map[string]environment.SerializableEnv
	CurrentEnvironment string
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (l *Manager) readState() error {
	if !fileExists(l.Options.StorageLocation) {
		return nil
	}
	f, err := ioutil.ReadFile(l.Options.StorageLocation)
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
	dat, err := json.Marshal(l.State) //, "", "    ")
	if err != nil {
		return err
	}

	ensureDir(l.Options.StorageLocation)

	return ioutil.WriteFile(l.Options.StorageLocation, dat, 0644)
}

func (l *Manager) ResetState() error {
	return os.Remove(l.Options.StorageLocation)
}

// Starts the local executor if it hasn't been started already
func (l *Manager) maybeStartLocalEnv() error {
	_, ok := l.State.Environments["local"]
	// At this point both Configuration and State are maps as they have been read from json

	// TODO: all this Executor management stuff and the state, is there a way to abstract it out?
	// Should it be in the executor interface itself?

	// As of now, we search the processes for the Wiz Executor PID, the State.Pid is not relevant
	// If we want to decode State.Pid we need mapstructure

	running := false
	var pid int32
	ps, err := process.Processes()
	if err != nil {
		return err
	}
	for _, p := range ps {
		n, err := p.Name()
		if err != nil {
			continue
		}

		if strings.Contains(n, "wiz") && p.Pid != int32(os.Getpid()) {
			pid = p.Pid
			log.Println("Found running Wiz process", n, pid)
			running = true

			// TODO: there could be multiple other wiz processes running, how to distinguish. RN we use the first one we find
			break
		}
	}

	log.Println(l.Options.RestartExecutor, running)

	//this section restarts the executor if its running
	if l.Options.RestartExecutor && running {
		log.Printf("Restarting executor process %d PID", pid)
		p, _ := process.NewProcess(pid)
		p.Kill()
		time.Sleep(200 * time.Millisecond)
		running = false
	}

	// If the environment isn't found or the process is not running, we create it
	if !ok || !running {
		config := local.Environment{Port: 8080}
		log.Println("Starting local environment with config", config)

		e := local.NewEnvironment()
		err := e.Configure(config)
		if err != nil {
			return err
		}

		data, err := e.StartExecutor("")
		if err != nil {
			return err
		}
		proc, ok := data.(*os.Process)
		if !ok {
			return fmt.Errorf("failed to convert executor response to Process, %v", data)
		}

		env := e.Describe()
		env.State = proc

		l.State.Environments["local"] = env
		l.State.CurrentEnvironment = "local"

		// Give the executor some time to start up and bind.
		time.Sleep(500 * time.Millisecond)
	}
	return nil
}

// RN it goes many Processors --> many downstream locations, etc
// TODO: think about centralizing this into a data structure that is easier to reason about

func setupProcessor(l Manager, pipeline tasks.Pipeline, node tasks.ProcessorNode) error {
	e := l.State.Environments[l.State.CurrentEnvironment]
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
	err := l.readState()
	if err != nil {
		return err
	}
	defer func() {
		err := l.writeState()
		if err != nil {
			log.Println(err)
		}
	}()

	log.Println("Creating pipeline:", p.Name)

	if _, ok := l.State.Pipelines[p.Name]; ok && !l.Options.OverwritePipelines {
		return fmt.Errorf("pipeline already exists")
	}
	localPipeline := &p
	localPipeline.UpdateFromSpec()
	err = localPipeline.CheckValidity()
	if err != nil {
		return err
	}

	log.Println("Pipeline", localPipeline.Name, "is valid, creating...")
	log.Println("Assigning runIDs to processors")

	localPipeline.AssignRunIDs(l.Options.PreserveRunIDs)
	localPipeline.UpdateInitialDataFlags()

	a := *localPipeline
	spew.Dump(a.Spec)

	l.State.Pipelines[p.Name] = a

	err = l.maybeStartLocalEnv()
	if err != nil {
		log.Println("failed to create env", err)
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
	e := manager.State.Environments[manager.State.CurrentEnvironment]
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

		chunkID := utils.GenID()

		log.Printf("Providing initial data to node %s (%s) with Chunk ID: %s", n.Name, n.Processor.ID, chunkID)

		return c.AddData(n.Processor.ID, n.RunID, api.Data{
			ChunkID:             chunkID,
			Format:              api.DataFormatRAW,
			Type:                api.DataTypeINPUT,
			State:               api.DataChunkStateWAITING,
			RawData:             p.Data,
			FilesystemReference: api.FilesystemReference{},
			AssociatedChunkID:   utils.GenID(),
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
