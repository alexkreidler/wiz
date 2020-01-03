package packages

import "github.com/alexkreidler/wiz/tasks"

// Package defines a serializable package
type Package struct {
	Name        string
	Description string
	Version     string
	Type        PackageType
	Source  	tasks.Pipeline
}