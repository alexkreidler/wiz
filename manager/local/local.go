package local

import (
	"encoding/json"
	"fmt"
	"github.com/alexkreidler/wiz/api"
	"github.com/alexkreidler/wiz/client"

	"github.com/alexkreidler/wiz/environment"
	"github.com/alexkreidler/wiz/environment/local"

	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/alexkreidler/wiz/tasks"
	"github.com/davecgh/go-spew/spew"
	"github.com/segmentio/ksuid"
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

func setupProcessor(l Manager, p tasks.ProcessorNode) error {
	e := l.Environments[l.CurrentEnvironment]
	if e.Host == "" {
		return fmt.Errorf("failed, invalid host")
	}

	id := p.Processor.Name
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

	runID := ksuid.New().String()
	log.Printf("Creating run %s for processor %s", runID, id)

	// POST /proc/id/run/id/config
	// Configure with Downstream True
	return c.Configure(id, runID, api.Configuration{
		ExpectedData: api.ExpectedData{
			NumChunks: 1,
		},
		ExecutorConfig: api.ExecutorConfig{
			SendDownstream: true,
			//DownstreamLocations: TODO
		},
		Processor: p.Processor.Configuration,
	})
}

func cpyM(m Manager) func(p tasks.ProcessorNode) error {
	return func(p tasks.ProcessorNode) error {
		return setupProcessor(m, p)
	}
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

	err = l.maybeStartLocalEnv()
	if err != nil {
		return err
	}

	err = localPipeline.Walk(cpyM(*l))
	if err != nil {
		log.Println(err)
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
