package httpjson

import (
	"encoding/json"
	"fmt"
	"github.com/alexkreidler/wiz/taguk"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"
)

type Server struct {
	Resources taguk.ResourceMap
	*mux.Router
	*http.Server
}

func NewServer() Server {
	return Server{
		make(taguk.ResourceMap),
		mux.NewRouter(),
		nil,
	}
}
func createPrimitiveObjects(t reflect.Type) reflect.Value {
	return reflect.Zero(t)
}

var specializedIndividualMethods = map[string]string{
	"GET": "Get",
	"POST": "Update",
	"DELETE": "Delete",
}
var specializedResourceMethods = map[string]string{
	"GET": "GetAll",
	"POST": "Create",
}

func (s *Server) Configure(rm taguk.ResourceMap) {
	s.Resources = rm


	for name, res := range rm {
		sr := s.Router.PathPrefix("/" + name).Subrouter()
		sr.HandleFunc("/{id:[0-9]+}/{action}", func(writer http.ResponseWriter, request *http.Request) {
			//res.Actions[]
			vars := mux.Vars(request)
			id := vars["id"]
			action := vars["action"]

			a, ok := res.Actions[action]
			if !ok {
				http.Error(writer, `{"error":"provided action is invalid"}`, 400)
				return
			}


			if n, err := strconv.ParseInt(id, 10, 64); err == nil {
				vs := a.Func.Call([]reflect.Value{createPrimitiveObjects(a.Base), reflect.ValueOf(n)})

				ret := vs[0]
				b, err := json.Marshal(ret.Interface())
				if err != nil {
					http.Error(writer, `{"error":"failed to marshal response"}`, 400)
					return
				}

				_, err = writer.Write(b)
				if err != nil {
					log.Println(err.Error())
				}
			} else {

				http.Error(writer, `{"error":"provided id is invalid"}`, 400)
				return
			}
		})
		for method, action := range specializedIndividualMethods {
			sr.NewRoute().Methods(method).Path("/{id:[0-9]+}").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				vars := mux.Vars(request)
				id := vars["id"]
				//fmt.Println(method, action)

				a, ok := res.Actions[taguk.FunctionNameToActionMap(action)]

				if !ok {
					http.Error(writer, `{"error":"provided action is invalid"}`, 400)
					return
				}
				spew.Dump(a)

				if n, err := strconv.ParseInt(id, 10, 64); err == nil {
					vs := a.Func.Call([]reflect.Value{createPrimitiveObjects(a.Base), reflect.ValueOf(n)})

					ret := vs[0]
					b, err := json.Marshal(ret.Interface())
					if err != nil {
						http.Error(writer, `{"error":"failed to marshal response"}`, 400)
						return
					}

					_, err = writer.Write(b)
					if err != nil {
						log.Println(err.Error())
					}
				} else {
					http.Error(writer, `{"error":"provided id is invalid"}`, 400)
					return
				}
			})
		}
		for method, action := range specializedResourceMethods {
			sr.NewRoute().Methods(method).Path("/").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				fmt.Println(method, action)
				fmt.Println("SUP")

				a, ok := res.Actions[taguk.FunctionNameToActionMap(action)]

				if !ok {
					http.Error(writer, `{"error":"provided action is invalid"}`, 400)
					return
				}
				spew.Dump(a)

				vs := a.Func.Call([]reflect.Value{createPrimitiveObjects(a.Base)})

				ret := vs[0]
				b, err := json.Marshal(ret.Interface())
				if err != nil {
					http.Error(writer, `{"error":"failed to marshal response"}`, 400)
					return
				}

				_, err = writer.Write(b)
				if err != nil {
					log.Println(err.Error())
				}
			})
		}

		//sr.HandleFunc("/{action}")
		//for _, a := range res.Actions {
		//	prefix := "/"
		//	if a.Individual {
		//		prefix += "{id}/"
		//	}
		//	switch a.Name {
		//	case "Create":
		//
		//	}
		//}
		//s.Router.HandleFunc("/" + name, func(writer http.ResponseWriter, request *http.Request) {
		//	res.
		//})
	}
}

func (s *Server) Serve() {
	p := os.Getenv("PORT")
	if p == "" {
		p = "8080"
	}

	hostname := "127.0.0.1"

	listener := hostname + ":" + p

	s.Server = &http.Server{
		Handler:      s.Router,
		Addr:         listener,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(s.Server.ListenAndServe())
}

