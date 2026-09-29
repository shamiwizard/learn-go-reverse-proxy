package main

import (
	"gopkg.in/yaml.v3"
	"os"
	"fmt"
	"maps"
	"slices"
	"strings"
)

var defaultRoute string

type stringSlice []string

var routes map[string]stringSlice

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

	var multi []string
	if err = value.Decode(&multi); err == nil {
		*r = multi
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

func matchPath(path string) string {
	if _, ok := routes[path]; ok {
		return path
	}

	rs := slices.SortedFunc(maps.Keys(routes), func(a, b string) int {
		return len(b) - len(a)
	})

	for _, route := range rs {
		if strings.HasPrefix(path, route) {
			return route
		}
	}

	return defaultRoute
}

func targets(path string) ([]string) {
	return routes[matchPath(path)]
}

