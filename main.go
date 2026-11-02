package main

import (
	"log"
	"net/http"
)

func main() {
	routes := setupRoutes()

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
