// Bxog is a simple and fast HTTP router for Go (HTTP request multiplexer).

package bxog

// Router
// API multiplexer, available methods:
//  New - create a new multiplexer
//  Add - add a rule specifying the handler (the default method - GET, ID - as a string to this rule)
//  Start - start the server indicating the listening port
//  Params - extract parameters from URL
//  Create - generate URL of the available options
//  Test - Start analogue (for testing only)
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	//"fmt"
	"log"
	"net/http"
	"time"
)

type Router struct {
	routes []*Route
	index  *Index
	url    string
}

func New() *Router {
	return &Router{}
}

func (r *Router) Add(url string, handler func(http.ResponseWriter, *http.Request, *Router)) *Route {
	if len(url) > HTTP_PATTERN_COUNT {
		panic("URL is too long")
		return nil
	} else {
		return r.newRoute(url, handler, HTTP_METHOD_DEFAULT)
	}
}

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

func (r *Router) Params(req *http.Request, id string) map[string]string {
	out := make(map[string]string)
	if c_route := r.index.index[r.index.genUint(id, 0)]; c_route != nil {
		query := c_route.genSplit(req.URL.Path[1:])
		for u := len(c_route.sections) - 1; u >= 0; u-- {
			if c_route.sections[u].type_sec == TYPE_ARG {
				out[c_route.sections[u].id] = query[u]
			}
		}
	} else {
		out["FUCK"] = "YOU"
	}
	return out
}

func (r *Router) Create(id string, param map[string]string) string {
	out := ""
	if route := r.index.index[r.index.genUint(id, 0)]; route != nil {
		for _, section := range route.sections {
			if section.type_sec == TYPE_STAT {
				out = out + DELIMITER_STRING + section.id
			} else if par, ok := param[section.id]; section.type_sec == TYPE_ARG && ok {
				out = out + DELIMITER_STRING + par
			}
		}
	}
	return out
}

func (r *Router) Test() {
	r.index = newIndex()
	r.index.compile(r.routes)
}
