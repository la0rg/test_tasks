package hash

import (
	"hash/crc32"
)

type Ring struct {
	head *Node
	left *Node
}

func (r *Ring) AddNode(name string) *Node {
	node := Node{
		name: name,
		hash: crc32.ChecksumIEEE([]byte(name)),
	}
	return r.addNode(&node)

}

func (r *Ring) Clear() {
	r.head = nil
	r.left = nil
}

func (r *Ring) addNode(node *Node) *Node {

	current := r.head

	if current == nil {
		r.head = node
		r.left = node
		return node
	}
	for {
		if node.hash < current.hash {
			if current.left == nil {
				current.left = node
				if current == r.left {
					r.left = node
				}
				return node
			}
			current = current.left
		} else {
			if current.right == nil {
				current.right = node
				return node
			}
			current = current.right
		}
	}
}

func (r *Ring) FindNode(key string) *Node {
	hash := crc32.ChecksumIEEE([]byte(key))
	return r.findNodeByHash(hash)
}

// The list of nodes that is responsible for storing a particular key.
// To account for node failures, preference list contains more
// than N nodes.
// Preference list for a key is constructed by skipping positions in the
// ring to ensure that the list contains only distinct physical nodes.
func (r *Ring) FindPreferenceList(key string, replication int) []*Node {
	node := r.FindNode(key)
	preferenceList := make([]*Node, 0)
	if node != nil {
		preferenceList = append(preferenceList, node)
		for len(preferenceList) < replication+1 {
			node := r.findSuccessorNode(node)
			preferenceList = append(preferenceList, node)
		}
	}
	return preferenceList
}

func (r *Ring) findNodeByHash(hash uint32) *Node {
	current := r.head

	// keys are connected to the next (clockwise) node on the ring

	// first step: search for a node based on hash
	for current != nil {
		if current.hash <= hash {
			if current.right != nil {
				current = current.right
			} else {
				break
			}
		} else {
			if current.left != nil {
				current = current.left
			} else {
				return current
			}
		}
	}

	// second step: search for the node successor
	return r.findSuccessorNode(current)
}

func (r *Ring) findSuccessorNode(current *Node) *Node {
	var head, successor *Node = r.head, nil
	if current != nil {

		// If right subtree of node is not nil, then succ lies in right subtree
		if current.right != nil {
			return minValue(current.right)
		}

		// Travel down the tree, if a node’s data is greater than root’s data then go right side,
		// otherwise go to left side.
		for head != nil {
			if current.hash < head.hash {
				successor = head
				head = head.left
			} else if current.hash > head.hash {
				head = head.right
			} else {
				break
			}
		}
	}
	if successor != nil {
		return successor
	}
	// for the most right node there is no successor -> close the ring
	return r.left
}

func minValue(node *Node) *Node {
	for node.left != nil {
		node = node.left
	}
	return node
}

type Node struct {
	name  string
	hash  uint32
	left  *Node
	right *Node
}
