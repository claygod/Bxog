// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

package bxog

// Server

import (
	//"fmt"
	"net/http"
)

// ServeHTTP looks for a matching route
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//r.index.find(req.URL.Path, req, r)
	//r.Default(w, req)
	//return
	if route := r.index.find(req.URL.Path, req, r); route != nil {
		//fmt.Println(route, " ??")
		//fmt.Println(" --> ", r)
		//fmt.Println(" --> ", r.current_route)
		route.handler(w, req, r)
		return
	} else {
		r.Default(w, req)
		return
	}
}

// Default Handler
func (r *Router) Default(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(404)
	http.Error(w, "Page not found", 404)
	return
}
