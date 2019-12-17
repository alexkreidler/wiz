package environment

import "fmt"

var Environments map[string]EnvironmentConstructor

type EnvironmentConstructor = func() Environment

func RegisterEnvironment(name string, constructor EnvironmentConstructor) error {
	if _, ok := Environments[name]; ok {
		return fmt.Errorf("failed to register environment with name %s, an environment was already registered", name)
	}
	Environments[name] = constructor
	return nil
}