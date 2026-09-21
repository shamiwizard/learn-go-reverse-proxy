package main

import (
	"gopkg.in/yaml.v3"
	"os"
	"fmt"
//	"io"
)

var routes Routes

type stringSlice []string

type Routes struct {
	Routes map[string]stringSlice `yaml:"routes"`
}

func init() {
	loadRoutes()
}

func (r *stringSlice) UnmarshalYAML(value *yaml.Node) error {
	var err error
	
	var single string
	if err = value.Decode(&single); err == nil {
		*r = []string{single}
		return err
	}

	var multy []string
	if err = value.Decode(&multy); err == nil {
		*r = multy
		return err
	}

	return fmt.Errorf("value must be string or []string")
}

func loadRoutes() {
	file, err := os.Open("./config/hosts.yml")

	if err != nil {
		panic("No hosts file is present")
	}

	defer file.Close()

	decoder := yaml.NewDecoder(file)

	err = decoder.Decode(&routes)

	if err != nil {
		fmt.Print("Something went wrong", err.Error())
	}
}

func target(path string) ([]string, error) {
	paths := routes.Routes[path]

	if len(paths) == 0 {
		return paths, fmt.Errorf("no paths are set for such target")
	}

	return paths, nil
}
