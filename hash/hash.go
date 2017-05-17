package hash

import (
	"fmt"
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

func (r *Ring) findNodeByHash(hash uint32) *Node {
	current := r.head

	// keys are connected to the next (clockwise) node on the ring

	// first step: search for a node based on hash
	for current != nil {
		if current.hash <= hash {
			if current.right != nil {
				fmt.Println("go right")
				current = current.right
			} else {
				break
			}
		} else {
			if current.left != nil {
				fmt.Println("go left")
				current = current.left
			} else {
				return current
			}
		}
	}

	var head, successor *Node = r.head, nil
	// second step: search for the node successor
	if current != nil {
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
	} else {
		// for the most right node there is no successor -> close the ring
		return r.left
	}
}

type Node struct {
	name  string
	hash  uint32
	left  *Node
	right *Node
}
