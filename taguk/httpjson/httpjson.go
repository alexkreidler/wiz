package httpjson

import (
	"encoding/json"
	"github.com/alexkreidler/wiz/taguk"
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
				writer.WriteHeader(400)
				_, _ = writer.Write([]byte(`{"error":"provided action is invalid"}`))
			}


			if n, err := strconv.ParseInt(id, 10, 64); err == nil {
				vs := a.Func.Call([]reflect.Value{createPrimitiveObjects(a.Base), reflect.ValueOf(n)})
				//vs
				//spew.Dump(vs[0])
				//json.M
				ret := vs[0]
				b, err := json.Marshal(ret.Interface())
				if err != nil {
					writer.WriteHeader(400)
					_, _ = writer.Write([]byte(err.Error()))
				}

				_, _ = writer.Write(b)
			} else {
				writer.WriteHeader(400)
				_, _ = writer.Write([]byte(`{"error":"provided id is invalid"}`))
			}

		})
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

