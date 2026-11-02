package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"maps"
	"os"
	"slices"
	"strings"
	"sync"
)

var configPath string = "./config/hosts.yml"

type stringSlice []string

type routes struct {
	routes       map[string]stringSlice
	defaultRoute string
}

type container struct {
	mu      sync.Mutex
	counter map[string]int
}

var con container = container{counter: make(map[string]int)}

func (r *stringSlice) UnmarshalYAML(node *yaml.Node) error {
	var err error

	var single string
	if err = node.Decode(&single); err == nil {
		*r = []string{single}
		return err
	}

	var multi []string
	if err = node.Decode(&multi); err == nil {
		*r = multi
		return err
	}

	return fmt.Errorf("value must be string or []string")
}

func (c *container) loadAndInc(route string, rs *routes) int {
	c.mu.Lock()

	defer c.mu.Unlock()

	if _, ok := c.counter[route]; !ok {
		c.counter[route] = 0
	}

	i := c.counter[route]
	c.counter[route]++

	if c.counter[route] >= len(rs.routes[route]) {
		c.counter[route] = 0
	}

	return i
}

func (rs *routes) load() {
	file, err := os.Open(configPath)

	if err != nil {
		panic("No hosts file is present")
	}

	defer file.Close()

	decoder := yaml.NewDecoder(file)

	err = decoder.Decode(&rs.routes)

	if err != nil {
		fmt.Print("Something went wrong", err.Error())
	}
}

func (rs *routes) setDefault() { // TODO: Change it to config method
	v, ok := rs.routes["default_route"]

	if !ok || v[0] == "" {
		panic(fmt.Sprintln("No default route is set.", "Set default_route in", configPath, "file."))
	}

	delete(rs.routes, "default_route")

	rs.defaultRoute = v[0]
}

func (rs *routes) normalize() {
	for k, v := range rs.routes {
		normalizedKey := normalizeRoute(k)

		if normalizedKey != k {
			rs.routes[normalizedKey] = v
			delete(rs.routes, k)
		}
	}
}

func normalizeRoute(route string) string {
	hasSuffix := strings.HasSuffix(route, "/")
	hasPrefix := strings.HasPrefix(route, "/")

	if hasPrefix && hasSuffix {
		return route
	}

	if !hasPrefix {
		route = "/" + route
	}

	if !hasSuffix {
		route += "/"
	}

	return route
}

func (rs *routes) match(path string) string {
	if _, ok := rs.routes[path]; ok {
		return path
	}

	r := slices.SortedFunc(maps.Keys(rs.routes), func(a, b string) int {
		return len(b) - len(a)
	})

	for _, route := range r {
		if strings.HasPrefix(path+"/", route) {
			return route
		}
	}

	return rs.defaultRoute
}

func (rs *routes) target(path string) string {
	route := rs.match(path)

	return rs.routes[route][con.loadAndInc(route, rs)]
}
