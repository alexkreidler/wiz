package main

import (
	"github.com/alexkreidler/deepcopy"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/src-d/go-git.v4"
	"log"
	"reflect"
)

func GetConfig() interface{} {
	return git.CloneOptions{}
}

func main() {
	baseConfig := GetConfig()
	opts := deepcopy.Copy(baseConfig, deepcopy.Options{ReturnPointer: true})

	//opts = opts.(git.CloneOptions)

	userConfig := map[string]interface{}{
		"Depth":        1,
		"SingleBranch": true,
	}

	log.Printf("%#+v", opts)
	v := reflect.ValueOf(opts).Elem() //.Elem()
	log.Println(v.CanSet(), v.Kind())

	err := mapstructure.Decode(userConfig, &opts)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%#+v", opts)
}
