// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

package bxog

// Index

import (
	"net/http"
	"strings"
)

// Router using the index selects the route
type index struct {
	tree  *node
	index map[typeHash]route
}

func newIndex() *index {
	return &index{
		index: make(map[typeHash]route),
	}
}

func (x *index) find(req *http.Request) *route {
	salt := x.genSalt(req.Method)
	cHashes := [HTTP_SECTION_COUNT]typeHash{}
	level := x.genUintSlice(req.URL.Path, salt, &cHashes)
	return x.findX(level, 0, x.tree, &cHashes)
}

func (x *index) findX(ln int, level int, tree2 *node, cHashes *[HTTP_SECTION_COUNT]typeHash) *route {
	if ln == level {
		return tree2.route
	} else if z1, ok := tree2.child[cHashes[level]]; ok {
		return x.findX(ln, level+1, z1, cHashes)
	} else if z2, ok := tree2.child[DELIMITER_UINT]; ok {
		return x.findX(ln, level+1, z2, cHashes)
	}
	return nil
}

func (x *index) getNode(arr map[string]*route) *node {
	out := newNode()
	childs := make(map[typeHash]map[string]*route)
	for url, r := range arr {
		url = strings.Trim(url, DELIMITER_STRING)
		if url == "" {
			out.route = r
			out.flag = true
			return out
		}
		arrStr := strings.Split(url, DELIMITER_STRING)
		if len(arrStr) > 1 {
			key := arrStr[0]
			salt := x.genSalt(r.method)
			hash := x.genUint(key, salt)
			if key[:1] == ":" {
				hash = DELIMITER_UINT
			}
			if _, ok := childs[hash]; !ok {
				childs[hash] = make(map[string]*route)
			}
			arrStr = arrStr[1:]
			url = strings.Join(arrStr, DELIMITER_STRING)
			childs[hash][url] = r
		} else if len(arrStr) == 1 {
			key := arrStr[0]

			salt := x.genSalt(r.method)
			hash := x.genUint(key, salt)
			if key[:1] == ":" {
				hash = DELIMITER_UINT
			}
			if _, ok := childs[hash]; !ok {
				childs[hash] = make(map[string]*route)
			}
			arrStr = make([]string, 0)
			url = ""
			childs[hash][url] = r
		}
	}
	for hash, v := range childs {
		n := x.getNode(v)
		out.child[hash] = n
	}
	return out
}

func (x *index) compile(routes []*route) {
	mapRoutes := make(map[string]*route)
	var core *route
	for _, r := range routes {
		x.index[x.genUint(r.id, 0)] = *r
		if r.url == "/" {
			core = r
		} else {
			mapRoutes[r.url] = r
		}
	}
	x.tree = x.getNode(mapRoutes)
	x.tree.route = core
}

func (x *index) genUintSlice(s string, salt typeHash, cHashes *[HTTP_SECTION_COUNT]typeHash) int {
	c := DELIMITER_BYTE
	na := 0
	length := len(s)
	if length == 1 {
		cHashes[0] = SLASH_HASH
		return 0
	}
	var total typeHash = salt
	for i := 1; i < length; i++ {
		if s[i] == c {
			cHashes[na] = total
			total = salt
			na++
			continue
		}
		total = total<<5 + typeHash(s[i])
	}
	cHashes[na] = total
	na++
	return na
}

func (x *index) genUint(s string, total typeHash) typeHash {
	length := len(s)
	for i := 0; i < length; i++ {
		total = total<<5 + typeHash(s[i])
	}
	return total
}

func (x *index) genSalt(s string) typeHash {
	return typeHash(s[0] + s[1])
}
