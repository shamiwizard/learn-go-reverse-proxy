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
			base, err := baseUrl(route)

			if err != nil {
				panic(err)
			}

			if _, ok := instances[base]; !ok {
				instances[base] = true
			}
		}
	}

	hc.instances = instances
}

func (hc *healthCheck) run(path string) {
	tick := time.Tick(5 * time.Second)

	for range tick {
		response, err := http.Get(healthUrl(path))

		if err != nil || response.StatusCode != http.StatusOK {
			//fmt.Println("Instnace", path, "disabled. Error:", err.Error())
			hc.instances[path] = false
			continue
		}

		hc.instances[path] = true
		//fmt.Println("Instance:", path, "are okay")
	}
}

func healthUrl(path string) string {
	return path + "/health"
}

func baseUrl(route string) (string, error) {
	u, err := url.Parse(route)

	if err != nil {
		return "", err
	}

	if u.Scheme == "" || u.Host == "" {
		return "", fmt.Errorf("Url does not have Scheme or Host")
	}

	return u.Scheme + "://" + u.Host, nil
}
