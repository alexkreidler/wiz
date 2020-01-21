package managers

// Manager represents a generic store that implements the appropriate CRUD operations on a given resource
// The manager may have specific reconciliation logic associated with the state of a given resource
type Manager interface {
	Create(r Resource) error
	Read(rid string) (Resource, error)
	Update(rid string, r Resource) error
	Delete(rid string) error
}

// A Resource represents any internal API object. It will be serialized to JSON for transport and internal API operations can be run on it once it is cast to the appropriate type.
type Resource struct {
	// Here we add additional options to change the behavior of the Go Restful implementation
	options interface{}

	data interface{}
}

// func (r Resource) Service() (restful.WebService)
