package main

import (
	"log"
	"net/http"
)

func main() {
	routes := setupRoutes()
	setupHealthCheck(routes)

	http.Handle("/", &requestResponder{routes: routes})

	log.Fatal(http.ListenAndServe(":1010", nil))
}

func setupRoutes() *routes {
	routes := new(routes)
	routes.load()
	routes.setDefault()
	routes.normalize()

	return routes
}

func setupHealthCheck(rs *routes) {
	check := healthCheck{rs: rs}
	check.start()
}
