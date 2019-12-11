package local

import (
	"fmt"
	"github.com/alexkreidler/wiz/executor"
	"github.com/alexkreidler/wiz/tasks"
)

type Manager struct {
	storageLocation string // the location that the local manager persists state to
	State
}

// State represents the manager state. It needs to be serializable to a file
type State struct {
	Pipelines    map[string]tasks.Pipeline
	Environments map[string]executor.SerializableEnv
}

func (l *Manager) CreatePipeline(p tasks.Pipeline, environment executor.Environment) error {
	if _, ok := l.State.Pipelines[p.Name]; ok {
		return fmt.Errorf("pipeline already exists")
	}
	l.State.Pipelines[p.Name] = p
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

func (l *Manager) CreateEnvironment(env executor.SerializableEnv) error {
	if _, ok := l.State.Environments[env.Name]; ok {
		return fmt.Errorf("environment already exists")
	}
	l.State.Environments[env.Name] = env
	return nil
}

func (l Manager) ReadEnvironment(name string) interface{} {
	env := l.State.Environments[name]
	return env.Configuration
}

func (l Manager) DeleteEnvironment(name string) error {
	if _, ok := l.State.Environments[name]; !ok {
		return fmt.Errorf("environment doesn't exist")
	}
	delete(l.State.Environments, name)
	return nil
}
