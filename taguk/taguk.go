package taguk

import "reflect"

// Server represents a taguk server
type Server struct {
	resources map[string]interface{}
}

// NewServer initializes a new server
func NewServer() Server {
	return Server{
		resources: make(map[string]interface{}),
	}
}

// AddResources adds several resources
func (s Server) AddResources(res ...interface{}) {
	for _, r := range res {
		s.AddResource(r)
	}
}

// AddResource adds the resource to the server, registering all actions for it as well
func (s *Server) AddResource(r interface{}) {
	// The following gets the type of the value underlying the interface{}
	v := reflect.ValueOf(r)
	if reflect.TypeOf(r).Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	
	rName := t.Name()
	println(rName)

	//println()

	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)

		println(m.Name)
	}
	
	s.resources[rName] = r
}