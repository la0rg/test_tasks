package hash

import "testing"

func TestAddNode(t *testing.T) {

	r := Ring{}
	r.AddNode("test 1")
}

func TestFindNode(t *testing.T) {

	//   5
	//  / \
	// 2   \
	//     10
	//    /  \
	//   6    15
	//    \
	//     8
	//
	r := Ring{}
	r.addNode(&Node{hash: 5})
	r.addNode(&Node{hash: 10})
	r.addNode(&Node{hash: 6})
	r.addNode(&Node{hash: 15})
	r.addNode(&Node{hash: 8})
	r.addNode(&Node{hash: 2})

	var ringTests = []struct {
		hash uint32
		node uint32
	}{
		{0, 2},
		{2, 5},
		{5, 6},
		{7, 8},
		{8, 10},
		{9, 10},
		{12, 15},
		{24, 2},
	}

	for _, ringTest := range ringTests {
		if r.findNodeByHash(ringTest.hash).hash != ringTest.node {
			t.Errorf("Expected node with hash %d", ringTest.node)
		}
	}
}

func TestFindNodeWithEmptyRing(t *testing.T) {
	r := Ring{}
	if r.FindNode("test") != nil {
		t.Error("Empty Ring should return nil")
	}
}
