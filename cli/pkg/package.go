package pkg

// Initial work on package data structures
// Should be in protobuf

// MetaInfo represents non-essential metadata
type MetaInfo struct {
	Author   string
	Repo     string
	Keywords string
}

// Dependency represents a Dependency of a package name and version
type Dependency struct {
	Name    string
	Version string
}

// Dependencies represent a list of package dependencies
type Dependencies []Dependency

//Package represents a package as defined in the docs
type Package struct {
	Name         string
	Version      string
	Type         string
	Meta         MetaInfo
	Dependencies Dependencies
}
