package test

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestMarshalBranch(t *testing.T) {
	b := Branch{}
	z := b.Get(12)

	byt, err := json.Marshal(z)
	Ok(t, err)
	fmt.Println(string(byt))
}
