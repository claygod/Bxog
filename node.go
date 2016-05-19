package bxog

// Node
// Nodes are stored in the index. Used for route search.
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

type Node struct {
	child map[type_hash]*Node
	route *Route
}

func newNode() *Node {
	new_node := &Node{}
	new_node.child = make(map[type_hash]*Node)
	return new_node
}
