package bxog

// Index
// Router using the index selects the route
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	//"fmt"
	"net/http"
)

type Index struct {
	tree  map[type_hash]*Node
	index map[type_hash]*Route
}

func newIndex() *Index {
	return &Index{
		tree:  make(map[type_hash]*Node),
		index: make(map[type_hash]*Route),
	}
}

func (x *Index) find(url string, req *http.Request, r *Router) *Route {
	salt := x.genSalt(req.Method)
	c_hashes := [HTTP_SECTION_COUNT]type_hash{}
	level := x.genUintSlice(url[1:], salt, &c_hashes)
	var c_node *Node

	//fmt.Println("== URL => ", url, " LEVEL => ", level)
	//fmt.Println("== HASH => ", c_hashes)
	//fmt.Println("== TREE => ", x.tree)

	if x.tree[c_hashes[0]] != nil {
		c_node = x.tree[c_hashes[0]]
	} else if x.tree[x.genUint(DELIMITER_STRING, salt)] != nil {
		c_node = x.tree[x.genUint(DELIMITER_STRING, salt)]
	} else {
		return nil
	}
	// slash "/"
	if level == 0 {
		if c_node.route != nil {
			return c_node.route
		} else {
			return nil
		}
	}

	for i := 0; i < level; i++ {
		if c_node.route == nil {
			if c_node.child[c_hashes[i]] != nil {
				c_node = c_node.child[c_hashes[i]]
			} else if c_node.child[DELIMITER_UINT] != nil {
				c_node = c_node.child[DELIMITER_UINT]
			} else {
				return nil
			}
		} else if i == level-1 {
			return c_node.route
		} else {
			return nil
		}
	}
	return nil
}

func (x *Index) compile(routes []*Route) {
	for _, route := range routes {
		salt := x.genSalt(route.method)
		x.index[x.genUint(route.id, 0)] = route
		length := len(route.sections)
		c_node := newNode()
		// slash
		if length == 0 {
			c_node.route = route
			x.tree[SLASH_HASH] = c_node
			continue
		}
		c_hash := x.genUint(route.sections[0].id, salt)
		if x.tree[c_hash] != nil {
			c_node = x.tree[c_hash]
		} else {
			switch route.sections[0].type_sec {
			case TYPE_STAT:
				x.tree[c_hash] = c_node
			case TYPE_ARG:
				x.tree[x.genUint(DELIMITER_STRING, salt)] = c_node
			}
		}
		// one
		/*
			if length == 1 {
				c_node.route = route
				continue
			}
		*/
		for i := 0; i < length; i++ {
			if i == length-1 {
				c_node.route = route
			} else {
				new_node := newNode()
				switch route.sections[i].type_sec {
				case TYPE_STAT:
					if i == 0 {
						c_node.child[x.genUint(route.sections[i].id, salt)] = new_node
					} else {
						c_node.child[x.genUint(route.sections[i].id, 0)] = new_node
					}
				case TYPE_ARG:
					if i == 0 {
						c_node.child[x.genUint(DELIMITER_STRING, salt)] = new_node
					} else {
						c_node.child[DELIMITER_UINT] = new_node
					}
				}
				c_node = new_node
			}
		}
	}
}

func (x *Index) genUintSlice(s string, total type_hash, c_hashes *[HTTP_SECTION_COUNT]type_hash) int {
	c := DELIMITER_BYTE
	na := 0
	length := len(s)
	if length == 0 {
		c_hashes[0] = SLASH_HASH
		return 0
	}

	for i := 0; i < length; i++ {
		if s[i] == c {
			c_hashes[na] = total
			total = 0
			na++
			continue
		}
		total = total<<5 + type_hash(s[i])
	}
	c_hashes[na] = total
	na++
	return na
}

func (x *Index) genUint(s string, total type_hash) type_hash {
	length := len(s)
	for i := 0; i < length; i++ {
		total = total<<5 + type_hash(s[i])
		//total = 31*total + type_hash(s[i])
	}
	return total
}

func (x *Index) genSalt(s string) type_hash {
	return type_hash(s[0] + s[1])
}
