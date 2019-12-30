package local

import (
	"encoding/json"
	"fmt"
	"github.com/alexkreidler/wiz/environment"
	"github.com/alexkreidler/wiz/environment/local"
	"github.com/alexkreidler/wiz/tasks"
	"github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"log"
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

func (l *Manager) readState() error {
	f, err := ioutil.ReadFile(l.storageLocation)
	if err != nil {
		return err
	}
	return json.Unmarshal(f, &l.State)
}

func (l *Manager) writeState() error {
	dat, err := json.Marshal(l.State)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(l.storageLocation, dat, 0644)
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

func (l *Manager) setupProcessor(p tasks.ProcessorNode) error {
	e := l.Environments[l.CurrentEnvironment]
	//endpoint := e.Endpoint
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
	log.Println("Pipeline", localPipeline.Name, "is valid, creating...")

	err = l.maybeStartLocalEnv()
	if err != nil {
		return err
	}
	spew.Dump(l)

	localPipeline.Walk(func(p tasks.ProcessorNode) error {
		return l.setupProcessor(p)
	})

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
