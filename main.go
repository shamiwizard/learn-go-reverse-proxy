package main

import (
	"log"
	"net/http"
)

func main() {
	routes := setupRoutes()
	healthCheck := setupHealthCheck(routes)

	http.Handle("/", &requestResponder{routes: routes, healthCheck: healthCheck})

	log.Fatal(http.ListenAndServe(":1010", nil))
}

func setupRoutes() *routes {
	routes := new(routes)
	routes.load()
	routes.setDefault()
	routes.normalize()

	return routes
}

func setupHealthCheck(rs *routes) *healthCheck {
	check := healthCheck{rs: rs}
	check.start()

	return &check
}
