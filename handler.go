package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"maps"
)

const ServerAddress = "http://localhost:2020"

type requestResponder struct{}

type flushWriter struct {
	w http.ResponseWriter
	f http.Flusher
}

func (fw flushWriter) Write(p []byte) (int, error) {
	n, err := fw.w.Write(p)

	if fw.f != nil {
		fw.f.Flush()
	}

	return n, err
}

func (rr *requestResponder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response, err := handleRequest(r)

	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	defer func() {
		response.Body.Close()
		logRequest(r)
	}()

	copyHeaders(response.Header, w.Header())

	w.WriteHeader(response.StatusCode)

	fw := flushWriter{w, w.(http.Flusher)}
	io.Copy(fw, response.Body)
}

func handleRequest(r *http.Request) (*http.Response, error) {
	newRequest, err := http.NewRequest(r.Method, fmt.Sprint(ServerAddress, r.URL.String()), r.Body)

	if err != nil {
		return nil, err
	}

	copyHeaders(r.Header, newRequest.Header)
	copyCookies(r, newRequest)
	setForwardHeaders(r, newRequest)

	response, err := http.DefaultClient.Do(newRequest)

	if err != nil {
		return nil, err
	}

	return response, err
}

func copyCookies(r, newRequest *http.Request) {
	cookies := r.Cookies()

	if len(cookies) == 0 {
		return
	}

	for _, cookie := range cookies {
		newCookie, err := http.ParseSetCookie(cookie.String())

		if err != nil {
			log.Printf("Cookie [%v: %v] is not added", cookie.Name, cookie.Value)
			continue
		}

		newRequest.AddCookie(newCookie)
	}
}

func copyHeaders(originalHeader, newHeader http.Header) {
	if originalHeader == nil {
		return
	}

	maps.Copy(newHeader, originalHeader)
}

func setForwardHeaders(originalRequest, newRequest *http.Request) {
	ip, _, err := net.SplitHostPort(originalRequest.RemoteAddr)

	if err == nil {
		newRequest.Header.Set("X-Forwarded-For", ip)
	}

	proto := "http"

	if originalRequest.TLS != nil {
		proto = "https"
	}

	newRequest.Header.Set("X-Forwarded-Proto", proto)

	newRequest.Header.Set("X-Forwarded-Host", originalRequest.Host)
}

