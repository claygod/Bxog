// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

package bxog

// route

import (
	"net/http"
	"strings"
)

// The route for URL
type route struct {
	id       string // added by the user
	method   string
	handler  func(http.ResponseWriter, *http.Request)
	sections []*section
}

func (r *Router) newRoute(url string, handler func(http.ResponseWriter, *http.Request), method string) *route {
	route := &route{
		url,
		method,
		handler,
		[]*section{},
	}
	route.setSections(url)
	r.routes = append(r.routes, route)
	return route
}

func (r *route) setSections(url string) {
	sec := r.parseUrl(url[1:])
	if len(sec) < HTTP_SECTION_COUNT {
		r.sections = sec
	} else {
		panic("Too many parameters!")
	}
}

func (r *route) Method(value string) *route {
	r.method = value
	return r
}

func (r *route) Id(value string) *route {
	r.id = value
	return r
}

func (r *route) parseUrl(url string) []*section {
	var arraySec []*section
	if len(url) == 0 {
		return []*section{}
	}
	result := r.genSplit(url)

	for _, value := range result {
		if strings.HasPrefix(value, ":") {
			arraySec = append(arraySec, newSection(value[1:], TYPE_ARG))
		} else {
			arraySec = append(arraySec, newSection(value, TYPE_STAT))
		}
	}
	return arraySec
}

func (r *route) genSplit(s string) []string {
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
