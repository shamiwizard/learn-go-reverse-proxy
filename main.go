package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/", new(requestResponder))

	log.Fatal(http.ListenAndServe(":1010", nil))
}
