package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type healthCheck struct {
	rs        *routes
	instances map[string]bool
}

func (hc *healthCheck) start() {
	hc.setInstances()
	fmt.Println(hc.instances)

	for path := range hc.instances {
		go hc.run(path)
	}
}

func (hc *healthCheck) setInstances() {
	instances := make(map[string]bool)

	for _, v := range hc.rs.routes {
		for _, route := range v {
			u, err := url.Parse(route)

			if err != nil {
				fmt.Println("Could not parse the route:", err)
				continue
			}

			if u.Scheme == "" || u.Host == "" {
				panic("One of the instances does not have hostname. Place check that instance has correct format")
			}

			hostname := u.Scheme + "://" + u.Host

			if _, ok := hc.instances[hostname]; !ok {
				instances[hostname] = true
			}
		}
	}

	hc.instances = instances
}

func (hc *healthCheck) run(path string) {
	tick := time.Tick(5 * time.Second)

	for range tick {
		response, err := http.Get(healthUrl(path))

		if err != nil {
			fmt.Println("Instnace", path, "disabled. Error:", err.Error())
			return
		}

		if response.StatusCode != http.StatusOK {
			fmt.Println("Instance", path, "is not working")
			return
		}

		fmt.Println("Instance:", path, "are okay")
	}
}

func healthUrl(path string) string {
	return path + "/health"
}
