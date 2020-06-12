// Copyright © 2016-2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

package bxog

// Index

import (
	//"fmt"
	"net/http"
	"strings"
)

// Router using the index selects the route
type index struct {
	tree       *node
	index      map[typeHash]route
	listShifts [DELIMITER_IN_LIST]int
	listRoutes []*route
	coreRoute  *route
	findMethod func(*http.Request) *route
}

func newIndex() *index {
	x := &index{
		index:      make(map[typeHash]route),
		listRoutes: []*route{&route{}}, // first route - dummy
	}

	x.findMethod = x.findTree

	return x
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

/*
When compiling the required data for two methods of routing search
(the second is left as an alternative)
*/
func (x *index) compile(routes []*route) {
	if len(routes) > MAX_ROUTES {
		panic("Too many routs, change the constant MAX_ROUTES in the configuration file.")
	}

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
	x.fillCore()
	x.fillNode(x.tree, 0)
}

func (x *index) fillCore() {
	x.coreRoute, x.tree.route = x.tree.route, x.coreRoute
}

func (x *index) fillNode(n *node, shiftLeft int) int {
	if shiftLeft > DELIMITER_IN_LIST-HTTP_PATTERN_COUNT {
		panic("Too many routs, change the constant MAX_ROUTES in the configuration file.")
	}

	shiftRigth := shiftLeft + 1 // zero - delim

	if n.route != nil {
		x.listShifts[shiftRigth] = -(int(len(x.listRoutes)))
		shiftRigth++
		x.listRoutes = append(x.listRoutes, n.route)

		return shiftRigth
	} else if countChild := len(n.child); countChild != 0 {
		shiftCur := shiftRigth
		shiftRigth += countChild * 2

		if _, ok := n.child[DELIMITER_UINT]; ok {
			x.listShifts[shiftCur] = DELIMITER_IN_LIST
			shiftCur++
			x.listShifts[shiftCur] = shiftRigth
			shiftRigth = x.fillNode(n.child[DELIMITER_UINT], shiftRigth)
			shiftCur++
		} else {
			for k, n2 := range n.child {
				x.listShifts[shiftCur] = int(k)
				shiftCur++
				x.listShifts[shiftCur] = shiftRigth
				shiftRigth = x.fillNode(n2, shiftRigth)
				shiftCur++
			}
		}

		return shiftRigth
	}

	return 0
}

func (x *index) findShift(req *http.Request) *route {
	cHashes := [HTTP_SECTION_COUNT]typeHash{}
	level := x.genUintSlice(req.URL.Path, x.genSalt(req.Method), &cHashes)

	if level == 0 && cHashes[0] == SLASH_HASH { //
		return x.coreRoute
	} else if level > 1 && cHashes[level-1] == SLASH_HASH { // 140
		return nil
	}

	curShift := 0

	for curLevel := 0; curLevel <= level; curLevel++ {
		shft := x.find2X(curLevel, curShift+1, cHashes[curLevel])

		switch {
		case shft < 0:
			shft = -shft

			return x.listRoutes[shft]

		case shft > 0:
			curShift = shft

		default:
			return nil
		}
	}

	return nil
}

func (x *index) find2X(curLevel int, curShift int, curHash typeHash) int {
	hash := int(curHash)

	for {
		csh := x.listShifts[curShift]

		if csh == hash {
			return x.listShifts[curShift+1]
		} else if csh == DELIMITER_IN_LIST {
			return x.listShifts[curShift+1]
		} else if csh < 0 {
			return csh
		} else if csh == 0 {
			return 0
		}

		curShift++
	}
}

func (x *index) findTree(req *http.Request) *route {
	cHashes := [HTTP_SECTION_COUNT]typeHash{}
	level := x.genUintSlice(req.URL.Path, x.genSalt(req.Method), &cHashes)

	if level == 0 && x.coreRoute != nil {
		return x.coreRoute
	}

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
