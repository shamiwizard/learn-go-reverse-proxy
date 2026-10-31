package main

import (
	"gopkg.in/yaml.v3"
	"os"
	"fmt"
	"maps"
	"slices"
	"strings"
	"sync"
)

var configPath string = "./config/hosts.yml"

var defaultRoute string

var routes map[string]stringSlice

type stringSlice []string

type container struct {
	mu sync.Mutex
	counter map[string]int
}

var con container = container{ counter: make(map[string]int) }

func init() {
	loadRoutes()
	setDefaultRoute()
	normalizeRoutes()
}

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

func (c *container) loadAndIncr(route string) int {
	c.mu.Lock()

	defer c.mu.Unlock()

	if _, ok := c.counter[route]; !ok {
		c.counter[route] = 0
	}

	i := c.counter[route]
	c.counter[route]++

	if c.counter[route] >= len(routes[route]) {
		c.counter[route] = 0
	}

	return i
}

func loadRoutes() {
	file, err := os.Open(configPath)

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

func setDefaultRoute() {
	v, ok := routes["default_route"]
	
	if !ok || v[0] == "" {
		panic(fmt.Sprintln("No default route is set.", "Set default_route in", configPath, "file."))
	} 

	delete(routes, "default_route")
	defaultRoute = v[0]
}

func normalizeRoutes() {
	for k, v := range routes {
		normalizedKey := normalizeRoute(k)

		if normalizedKey != k {
			routes[normalizedKey] = v
			delete(routes, k)
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

func matchPath(path string) string {
	if _, ok := routes[path]; ok {
		return path
	}

	rs := slices.SortedFunc(maps.Keys(routes), func(a, b string) int {
		return len(b) - len(a)
	})

	for _, route := range rs {
		if strings.HasPrefix(path + "/", route) {
			return route
		}
	}

	return defaultRoute
}

func target(path string) string {
	route := matchPath(path)

	return routes[route][con.loadAndIncr(route)]
}

