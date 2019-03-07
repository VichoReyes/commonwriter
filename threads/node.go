package threads

import "sync"

// Node is an immutable snapshot of the story, like a commit
// its Children are the story versions based on it
// its content can be rendered using the String method
type Node struct {
	content string // TODO turn into diff
	// children map[int]*Node
	children []*Node
	Title    string
	Authors  map[string]bool
	sync.Mutex
}

// Content returns the whole story
// TODO change signature to HTML to allow some markup
func (n *Node) Content() string {
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

// Append makes a node n get a child with content, author and title
func (n *Node) Append(content, author, title string) {
	var new Node
	new.content = content
	new.Authors = cloneSet(n.Authors, author)
	new.Title = title
	n.Lock()
	n.children = append(n.children, &new)
	n.Unlock()
}

func cloneSet(old map[string]bool, elem string) map[string]bool {
	if old[elem] { // old already has elem
		return old
	}
	ret := make(map[string]bool)
	for k := range old {
		ret[k] = true
	}
	ret[elem] = true
	return ret
}
