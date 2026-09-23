package main

import (
	"gopkg.in/yaml.v3"
	"os"
	"fmt"
	"maps"
	"slices"
	"regexp"
)

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

func matchRoute(path string) string {
	if _, ok := routes[path]; ok {
		return path
	}

	keys := slices.Sorted(maps.Keys(routes))

	result := keys[0]
	var re *regexp.Regexp
	var err error

	for _, v := range keys {
		re, err = regexp.Compile(fmt.Sprint("(^", regexp.QuoteMeta(v), ")"))

		if err != nil {
			break
		}

		if re.MatchString(path) {
			result = v
		}
	}

	return result
}

func targets(path string) ([]string) {
	return routes[matchRoute(path)]
}


