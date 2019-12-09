package tasks

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"testing"
)

func TestTasksDeserialization(t *testing.T) {
	bytes, err := ioutil.ReadFile("./test.yaml")
	if err != nil {
		t.Fatal(err)
	}

	p := Pipeline{}

	err = yaml.Unmarshal(bytes, &p)
	if err != nil {
		t.Fatal(err)
	}

	//fmt.Printf("%#v", p)
	spew.Dump(p)


	out, err := yaml.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(out))
}

func TestTasksGraphCreation(t *testing.T) {
	bytes, err := ioutil.ReadFile("./test.yaml")
	if err != nil {
		t.Fatal(err)
	}

	p := NewPipeline("test")

	err = yaml.Unmarshal(bytes, &p)
	if err != nil {
		t.Fatal(err)
	}

	p.UpdateFromSpec()

	//fmt.Printf("%#v", p)
	spew.Dump(p)
}

func TestTasksGraphVisualization(t *testing.T) {
	bytes, err := ioutil.ReadFile("./test.yaml")
	if err != nil {
		t.Fatal(err)
	}

	p := NewPipeline("test")

	err = yaml.Unmarshal(bytes, &p)
	if err != nil {
		t.Fatal(err)
	}

	p.UpdateFromSpec()

	//fmt.Printf("%#v", p)
	//spew.Dump(p)
	out, err := dot.Marshal(p.g, "test graph", "", "")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(out))
}

func TestTasksSerialization(t *testing.T) {
	p := Pipeline{
		Name:     "test",
		g:        nil,
		rootNode: nil,
		Spec: PipelineSpec{
			Sequential: Sequential{
				ProcessorNode{
					id:   0,
					Name: "",
					Processor: Processor{
						Type:          "",
						Version:       "",
						Configuration: nil,
					},
					Children: Children{
						Sequential: nil,
						Parallel:   nil,
					},
				},
			},
		},
	}

	out, err := yaml.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(out))
}
