package get

import (
	"sync"
	"testing"

	"gotest.tools/assert"
)

func TestBasicGet(t *testing.T) {
	p := GetProcessor{}
	err := p.Configure(&GoGetConfig{Source:""})
	assert.NilError(t, err)

	t.Log("testing get")

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		t.Log("here")
		for state := range p.State() {
			t.Log("state change", state)
		}
		wg.Done()
	}()


	sources := []string{"https://google.com"}//,"github.com/hashicorp/go-getter.git"}
	for _, source := range sources {
		t.Log("Testing", source)
		data := map[string]interface{}{"Source": source}
		t.Log(data)
		//wg.Add(1)
		go func() {
			wg.Add(1)
			p.Run(data)
			wg.Done()
		}()
	}

	wg.Wait()
	//p.dir
	t.Log(p.dir)
}