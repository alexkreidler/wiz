package local

import (
	"fmt"
	"testing"
)

func TestStartLocalExecutor(t *testing.T) {
	e := Environment{Port: 9003}

	err := e.StartExecutor("random")
	if err != nil {
		t.Fatal("failed to start executor")
	}
}

func TestGetEnvironmentName(t *testing.T) {
	e := Environment{Port: 9003}
	fmt.Println(e.Name())
}