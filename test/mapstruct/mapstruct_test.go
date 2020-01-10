package mapstruct

import (
	jsonscrape "github.com/alexkreidler/jsonscrape/lib"
	"github.com/mitchellh/mapstructure"
	"testing"
)

// Note that the mapstructure tags defined in the struct type
// can indicate which fields the values are mapped to.
type Person struct {
	Name string `mapstructure:"person_name"`
	Age  int    `mapstructure:"person_age"`
}

func TestMapStruct(t *testing.T) {

	input := map[string]interface{}{
		"person_name": "Mitchell",
		"person_age":  91,
	}

	var result Person
	err := mapstructure.Decode(input, &result)
	if err != nil {
		panic(err)
	}

	t.Logf("%#v", result)

}

func TestDecodeConfig(t *testing.T) {
	input := map[string]interface{}{
		"config": map[string]interface{}{
			"Sites": []interface{}{
				"rcp.com", "latest_polls",
			},
		},
	}

	var config jsonscrape.Config

	err := mapstructure.Decode(input, &config)
	if err != nil {
		panic(err)
	}

	t.Logf("%#v", config)

	if len(config.GeneralConfig.Sites) != 2 {
		t.Error("decoded struct doesn't match input")
	}
}
