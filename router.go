// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

package bxog

// Router

import (
	"log"
	"net/http"
	"time"
)

// Router Bxog is a simple and fast HTTP router for Go (HTTP request multiplexer).
type Router struct {
	routes []*route
	index  *index
	url    string
}

// New - create a new multiplexer
func New() *Router {
	return &Router{}
}

// Add - add a rule specifying the handler (the default method - GET, ID - as a string to this rule)
func (r *Router) Add(url string, handler func(http.ResponseWriter, *http.Request)) *route {
	if len(url) > HTTP_PATTERN_COUNT {
		panic("URL is too long")
	} else {
		return r.newRoute(url, handler, HTTP_METHOD_DEFAULT)
	}
}

// Start - start the server indicating the listening port
func (r *Router) Start(port string) {
	r.index = newIndex()
	r.index.compile(r.routes)
	s := &http.Server{
		Addr:           port,
		Handler:        nil,
		ReadTimeout:    READ_TIME_OUT * time.Second,
		WriteTimeout:   WRITE_TIME_OUT * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	http.Handle(DELIMITER_STRING, r)
	http.Handle(FILE_PREF, http.StripPrefix(FILE_PREF, http.FileServer(http.Dir(FILE_PATH))))
	log.Fatal(s.ListenAndServe())
}

// Test - Start analogue (for testing only)
func (r *Router) Test() {
	r.index = newIndex()
	r.index.compile(r.routes)
}
