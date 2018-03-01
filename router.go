// Copyright © 2016-2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

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
func (r *Router) Add(url string, handler func(http.ResponseWriter, *http.Request, *Router)) *route {
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

// Params - extract parameters from URL
func (r *Router) Params(req *http.Request, id string) map[string]string {
	out := make(map[string]string)
	if cRoute, ok := r.index.index[r.index.genUint(id, 0)]; ok {
		query := cRoute.genSplit(req.URL.Path[1:])
		for u := len(cRoute.sections) - 1; u >= 0; u-- {
			if cRoute.sections[u].typeSec == TYPE_ARG {
				out[cRoute.sections[u].id] = query[u]
			}
		}
	}
	return out
}

// Create - generate URL of the available options
func (r *Router) Create(id string, param map[string]string) string {
	out := ""
	if route, ok := r.index.index[r.index.genUint(id, 0)]; ok {
		for _, section := range route.sections {
			if section.typeSec == TYPE_STAT {
				out = out + DELIMITER_STRING + section.id
			} else if par, ok := param[section.id]; section.typeSec == TYPE_ARG && ok {
				out = out + DELIMITER_STRING + par
			}
		}
	}
	return out
}

// Test - Start analogue (for testing only)
func (r *Router) Test() {
	r.index = newIndex()
	r.index.compile(r.routes)
}
