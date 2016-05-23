// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

package bxog

// Node

// Nodes are stored in the index. Used for route search.
type node struct {
	child map[type_hash]*node
	route *route
}

func newNode() *node {
	new_node := &node{}
	new_node.child = make(map[type_hash]*node)
	return new_node
}
