package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

//CONST is a global var
var CONST = "Hello"

// ReadSpecFile reads and prints a spec file
func ReadSpecFile(file string) {
	fmt.Println("Read Spec")
	dat, err := ioutil.ReadFile(file)
	check(err)
	fmt.Print(string(dat))
	var i interface{}
	err = json.Unmarshal(dat, &i)
	fmt.Println(i)
	check(err)

	m := i.(map[string]interface{})

	for k, v := range m {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "is string", vv)
		case int:
			fmt.Println(k, "is int", vv)
		case []interface{}:
			fmt.Println(k, "is an array:")
			for i, u := range vv {
				fmt.Println(i, u)
			}
		case interface{}:
			fmt.Println(k, "is a map", vv)
		default:
			fmt.Println(k, "is of a type I don't know how to handle")
		}
	}
}
