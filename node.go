// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

package bxog

// Node

// Nodes are stored in the index. Used for route search.
type node struct {
	child map[typeHash]*node
	route *route
}

func newNode() *node {
	nNode := &node{}
	nNode.child = make(map[typeHash]*node)
	return nNode
}
