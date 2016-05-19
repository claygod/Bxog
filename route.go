package bxog

// Route
// The route for URL
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	//"fmt"
	"net/http"
	"strings"
)

type Route struct {
	id       string // added by the user
	method   string
	handler  func(http.ResponseWriter, *http.Request, *Router)
	sections []*Section
}

func (r *Router) newRoute(url string, handler func(http.ResponseWriter, *http.Request, *Router), method string) *Route {
	route := &Route{
		url,
		method,
		handler,
		[]*Section{},
	}
	//route.id = url
	//route.handler = handler
	//route.method = method
	route.setSections(url)
	r.routes = append(r.routes, route)
	return route
}

func (r *Route) setSections(url string) *Route {
	sec := r.parseUrl(url[1:])
	if len(sec) < HTTP_SECTION_COUNT {
		r.sections = sec
		return r
	} else {
		panic("Too many parameters!")
		return nil
	}
}

func (r *Route) Method(value string) *Route {
	r.method = value
	return r
}

func (r *Route) Id(value string) *Route {
	r.id = value
	return r
}

func (r *Route) parseUrl(url string) []*Section {
	var array_sec []*Section
	if len(url) == 0 {
		return []*Section{}
	}
	result := r.genSplit(url)

	for _, value := range result {
		if strings.HasPrefix(value, ":") {
			array_sec = append(array_sec, newSection(value[1:], TYPE_ARG))
		} else {
			array_sec = append(array_sec, newSection(value, TYPE_STAT))
		}
	}
	return array_sec
}

func (r *Route) genSplit(s string) []string {
	n := 1
	c := DELIMITER_BYTE
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			n++
		}
	}
	out := make([]string, n)
	count := 0
	begin := 0
	length := len(s) - 1
	for i := 0; i <= length; i++ {
		if s[i] == c {
			out[count] = s[begin:i]
			count++
			begin = i + 1
		}
	}
	out[count] = s[begin : length+1]
	return out
}
