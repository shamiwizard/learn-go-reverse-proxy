package main

import (
	"log"
	"net/http"
)

func logRequest(r *http.Request) {
	path := r.URL.String()
	method := r.Method

	log.Printf("--- %v %v ---", method, path)
}
