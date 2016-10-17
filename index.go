// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

package bxog

// Index

import (
	"net/http"
)

// Router using the index selects the route
type index struct {
	tree  map[typeHash]*node
	index map[typeHash]*route
}

func newIndex() *index {
	return &index{
		tree:  make(map[typeHash]*node),
		index: make(map[typeHash]*route),
	}
}

func (x *index) find(url string, req *http.Request) *route {
	salt := x.genSalt(req.Method)
	cHashes := [HTTP_SECTION_COUNT]typeHash{}
	level := x.genUintSlice(url[1:], salt, &cHashes)
	var cNode *node

	if x.tree[cHashes[0]] != nil {
		cNode = x.tree[cHashes[0]]
	} else if x.tree[x.genUint(DELIMITER_STRING, salt)] != nil {
		cNode = x.tree[x.genUint(DELIMITER_STRING, salt)]
	} else {
		return nil
	}
	// slash
	if level == 0 {
		return cNode.route
	}

	for i := 0; i < level; i++ {
		if cNode.route == nil {
			if cNode.child[cHashes[i]] != nil {
				cNode = cNode.child[cHashes[i]]
			} else if cNode.child[DELIMITER_UINT] != nil {
				cNode = cNode.child[DELIMITER_UINT]
			} else {
				return nil
			}
		} else if i == level-1 {
			return cNode.route
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
		cNode := newNode()
		// slash
		if length == 0 {
			cNode.route = route
			x.tree[SLASH_HASH] = cNode
			continue
		}
		cHash := x.genUint(route.sections[0].id, salt)
		if x.tree[cHash] != nil {
			cNode = x.tree[cHash]
		} else {
			switch route.sections[0].typeSec {
			case TYPE_STAT:
				x.tree[cHash] = cNode
			case TYPE_ARG:
				x.tree[x.genUint(DELIMITER_STRING, salt)] = cNode
			}
		}
		for i := 0; i < length; i++ {
			if i == length-1 {
				cNode.route = route
			} else {
				nNode := newNode()
				switch route.sections[i].typeSec {
				case TYPE_STAT:
					if i == 0 {
						cNode.child[x.genUint(route.sections[i].id, salt)] = nNode
					} else {
						cNode.child[x.genUint(route.sections[i].id, 0)] = nNode
					}
				case TYPE_ARG:
					if i == 0 {
						cNode.child[x.genUint(DELIMITER_STRING, salt)] = nNode
					} else {
						cNode.child[DELIMITER_UINT] = nNode
					}
				}
				cNode = nNode
			}
		}
	}
}

func (x *index) genUintSlice(s string, total typeHash, cHashes *[HTTP_SECTION_COUNT]typeHash) int {
	// fmt.Print("\n -- salt: ", typeHash)
	c := DELIMITER_BYTE
	na := 0
	length := len(s)
	if length == 0 {
		cHashes[0] = SLASH_HASH
		return 0
	}

	for i := 0; i < length; i++ {
		if s[i] == c {
			cHashes[na] = total
			total = 0
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
