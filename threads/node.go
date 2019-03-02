package threads

import "sync"

// Node is an immutable snapshot of the story, like a commit
// its Children are the story versions based on it
// its content can be rendered using the String method
type Node struct {
	content string // TODO turn into diff
	// children map[int]*Node
	children []*Node
	sync.Mutex
}

func (n *Node) String() string {
	n.Lock()
	defer n.Unlock()
	return n.content
}

// Children returns the list of Nodes based on n
// no guarantees are made about their order (for now)
func (n *Node) Children() []*Node {
	/*
		chil := make([]*Node, len(n.children))
		n.Lock()
		for _, v := range n.children {
			chil = append(chil, v)
		}
		n.Unlock()
		return chil
	*/
	return n.children
}

// Append makes a node n get a child with content succesor
func (n *Node) Append(succesor string) {
	var new Node
	new.content = succesor
	n.Lock()
	n.children = append(n.children, &new)
	n.Unlock()
}
