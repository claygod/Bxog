// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

package bxog

// Index

import (
	"net/http"
)

// Router using the index selects the route
type index struct {
	tree  map[type_hash]*node
	index map[type_hash]*route
}

func newIndex() *index {
	return &index{
		tree:  make(map[type_hash]*node),
		index: make(map[type_hash]*route),
	}
}

func (x *index) find(url string, req *http.Request, r *Router) *route {
	salt := x.genSalt(req.Method)
	c_hashes := [HTTP_SECTION_COUNT]type_hash{}
	level := x.genUintSlice(url[1:], salt, &c_hashes)
	var c_node *node

	if x.tree[c_hashes[0]] != nil {
		c_node = x.tree[c_hashes[0]]
	} else if x.tree[x.genUint(DELIMITER_STRING, salt)] != nil {
		c_node = x.tree[x.genUint(DELIMITER_STRING, salt)]
	} else {
		return nil
	}
	// slash
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

func (x *index) compile(routes []*route) {
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

func (x *index) genUintSlice(s string, total type_hash, c_hashes *[HTTP_SECTION_COUNT]type_hash) int {
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

func (x *index) genUint(s string, total type_hash) type_hash {
	length := len(s)
	for i := 0; i < length; i++ {
		total = total<<5 + type_hash(s[i])
	}
	return total
}

func (x *index) genSalt(s string) type_hash {
	return type_hash(s[0] + s[1])
}
