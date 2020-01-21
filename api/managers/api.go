package managers
import (
	"encoding/json"
	restful "github.com/emicklei/go-restful/v3"
	"github.com/iancoleman/strcase"
	"io/ioutil"
)
// Manager represents a generic store that implements the appropriate CRUD operations on a given resource
// The manager may have specific reconciliation logic associated with the state of a given resource
type Manager interface {
	Create(r Resource) error
	Read(rid string) (Resource, error)
	Update(rid string, r Resource) error
	Delete(rid string) error
}

// A Resource represents any internal API object. It will be serialized to JSON for transport and internal API operations can be run on it once it is cast to the appropriate type.
// Here we add additional options to change the behavior of the Go Restful implementation
// Also, all resources will be merged with the default resource to set defaults
type Resource struct {
	// The name automatically generates the HTTP route
	Name string
	Consumes []string
	Produces []string

	// Internal is the internal struct API representation
	Internal IResource

	//data interface{}
}

type IResource interface {
	Create(r interface{}) error
	Read(rid string) (interface{}, error)
	Update(rid string, r interface{}) error
	Delete(rid string) error
}



func (r Resource) Service() *restful.WebService {
	ws := new(restful.WebService)
	ws.
		Path(strcase.ToKebab(r.Name)).
		Consumes(r.Consumes...). // by default these formats are JSON
		Produces(r.Produces...)

	ws.GET("/").To(func(request *restful.Request, response *restful.Response) {
		body := request.Request.Body

		buf, err := ioutil.ReadAll(body)
		if err != nil {
			_ = response.WriteErrorString(500, err.Error())
		}

		l := r.Internal
		err = json.Unmarshal(buf, &l)
		if err != nil {
			_ = response.WriteErrorString(500, err.Error())
		}
		//r.Internal.Create()
	}).
		Doc("get a " + r.Name)
		//Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
		//Writes(User{}))

	return ws
}
