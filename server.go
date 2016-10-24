// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

package bxog

// Server

import (
	"context"
	"net/http"
)

// ServeHTTP looks for a matching route
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if route := r.index.find(req.URL.Path, req); route != nil {
		query := route.genSplit(req.URL.Path[1:])
		ctx := req.Context()
		for u := len(route.sections) - 1; u >= 0; u-- {
			if route.sections[u].typeSec == TYPE_ARG {
				ctx = context.WithValue(ctx, route.sections[u].id, query[u])
			}
		}
		req = req.WithContext(ctx)
		route.handler(w, req)
	} else {
		r.Default(w, req)
	}
}

// Default Handler
func (r *Router) Default(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(404)
	http.Error(w, "Page not found", 404)
}
