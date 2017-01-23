// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

package bxog

// Server

import (
	"net/http"
)

// ServeHTTP looks for a matching route
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if route := r.index.find(req); route != nil {
		route.handler(w, req, r)
	} else {
		r.Default(w, req)
	}
}

// Default Handler
func (r *Router) Default(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(404)
	//http.Error(w, "Page not found", 404)
}
