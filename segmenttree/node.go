package segmenttree

import "math"

type Node struct {
	/*
			Size of a node is
			keys: 4+4+4+ byte (reference, length, cap)  + b-1 * 4 byte
			values: assume float32 --> 12 byte + b * 4 byte
			children: 12 byte + l * 4 byte
			parent: 4 byte
			tree: 4 byte
			isLeaf: 1 byte

		This gives a total of  3*12 +(2*b-1) *4 + l * 4 + 8 +1 = 45 + (2*b-1) *4 + l * 4

		Go determines the memory allocation for our structs, it will pad bytes to make sure the final memory
		footprint is a multiple of 8 bytes
		let's assume the default page size of linux with 4 kb = 4096 bytes

		For simplicity set b = l, then we get a maximal branching factor b and a maximal leaf capacity l
		of b=l=337 to fit into one disk page of 4 kb.
	*/
	keys     []uint32
	values   []Addable
	children []*Node
	parent   *Node
	tree     *SegmentTreeImpl
	isLeaf   bool
}

func (node *Node) findIntervalIndex(instant uint32) uint32 {
	var intervalIndex uint32 = 0

	for intervalIndex < node.size() {
		if instant >= node.keys[intervalIndex] {
			intervalIndex++
		} else {
			break
		}
	}

	return intervalIndex
}

func (node *Node) findChildIndex(childToFind *Node) uint32 {
	for index, child := range node.children {
		if child == childToFind {
			return uint32(index)
		}
	}

	panic("Child is not a child of parent")
}

func (node *Node) getIntervalStart(index uint32) uint32 {
	if index == 0 {
		if node.parent == nil {
			return 0 // min of Uint32
		} else {
			childIndex := node.parent.findChildIndex(node)
			return node.parent.getIntervalStart(childIndex)
		}
	} else {
		return node.keys[index-1]
	}
}

func (node *Node) getIntervalEnd(index uint32) uint32 {
	if index >= node.size() {
		if node.parent == nil {
			return math.MaxUint32
		} else {
			childIndex := node.parent.findChildIndex(node)
			return node.parent.getIntervalEnd(childIndex)
		}
	} else {
		return node.keys[index]
	}
}

func (node *Node) getIntervals() []Interval {
	var intervals []Interval = make([]Interval, node.size()+1)

	var i uint32 = 0
	for ; i <= node.size(); i++ {
		start := node.getIntervalStart(i)
		end := node.getIntervalEnd(i)

		intervals[i] = NewInterval(start, end)
	}

	return intervals
}

func (node *Node) insert(intervalIndex int, tupleToInsert ValueIntervalTuple) int {
	nodeIntervalStart := node.getIntervalStart(uint32(intervalIndex))
	nodeIntervalEnd := node.getIntervalEnd(uint32(intervalIndex))

	if nodeIntervalStart < tupleToInsert.interval.start && nodeIntervalEnd > tupleToInsert.interval.end {
		// Case 1
		// Node interval:   |---------|
		// Insert interval:   |-----|
		// The cases where they are equal or where the node interval is a subset
		// of the insert interval is handled at tree level.

		// Make space for two additional intervals
		node.keys = append(node.keys, 0, 0)
		node.values = append(node.values, Float(0), Float(0))

		// Shift keys and insert
		for i := len(node.keys) - 1; i > intervalIndex+1; i-- {
			node.keys[i] = node.keys[i-2]
		}
		node.keys[intervalIndex] = tupleToInsert.interval.start
		node.keys[intervalIndex+1] = tupleToInsert.interval.end

		// Shift values and insert
		for i := len(node.keys); i > intervalIndex+2; i-- {
			node.values[i] = node.values[i-2]
		}

		node.values[intervalIndex+2] = node.values[intervalIndex] // Original value
		node.values[intervalIndex+1] = node.tree.aggregate.operation(node.values[intervalIndex], tupleToInsert.value)

		return 2 // one interval got split into three (= 2 new)
	} else if nodeIntervalStart < tupleToInsert.interval.start && nodeIntervalEnd <= tupleToInsert.interval.end {
		// Case 2
		// Node interval:      |--------|
		// Insert interval:      |------|~~|

		// Add space for a new interval
		node.keys = append(node.keys, 0)
		node.values = append(node.values, Float(0))

		// Shift keys and insert
		for i := len(node.keys) - 1; i > intervalIndex; i-- {
			node.keys[i] = node.keys[i-1]
		}

		node.keys[intervalIndex] = tupleToInsert.interval.start

		// Shift values and insert
		for i := len(node.keys); i > intervalIndex; i-- {
			node.values[i] = node.values[i-1]
		}

		node.values[intervalIndex+1] = node.tree.aggregate.operation(node.values[intervalIndex+1], tupleToInsert.value)

		return 1 // one interval got split into two (= 1 new)
	} else if nodeIntervalStart >= tupleToInsert.interval.start && nodeIntervalEnd > tupleToInsert.interval.end {
		// Case 3
		// Node interval:      |--------|
		// Insert interval: |~~|------|

		// Add space for a new interval
		node.keys = append(node.keys, 0)
		node.values = append(node.values, Float(0))

		// Shift keys and insert
		for i := len(node.keys) - 1; i > intervalIndex; i-- {
			node.keys[i] = node.keys[i-1]
		}

		node.keys[intervalIndex] = tupleToInsert.interval.end

		// Shift values and insert
		for i := len(node.keys); i > intervalIndex; i-- {
			node.values[i] = node.values[i-1]
		}

		node.values[intervalIndex] = node.tree.aggregate.operation(node.values[intervalIndex], tupleToInsert.value)

		return 1 // one interval got split into two (= 1 new)
	} else {
		// All other cases are handled at tree level.
		// Thus, we should not get here.
		panic("should not get here")
	}
}

func (node *Node) split() *Node {
	if node.size() < 1 {
		// Let's add an invariant to get rid of ugly edge cases, which are irrelevant in practice!
		panic("A Node of size < 2 can not be split.")
	}
	if node.size()+1 <= node.tree.branchingFactor || node.size()+1 > node.tree.branchingFactor+2 {
		// Let's add an invariant to get rid of ugly edge cases, which are irrelevant in practice!
		panic("A Node to split has to have a size + 1  of l + 1 or l+2 in case it is a leaf and b+1 in case it is not a leaf, where l = b = branchingFactor.")
	}
	var parent *Node
	n := node.size() + 1
	half_n := int32(math.Ceil(float64(n) / float64(2)))

	// N1 contains 1 ... n/2-1 instances and corresponding pointers if not a leaf child
	n1 := &Node{
		keys:     node.keys[:half_n-1],
		values:   node.values[:half_n],
		children: nil,
		parent:   node.parent,
		tree:     node.tree,
		isLeaf:   node.isLeaf,
	}

	// N2 contains n/2 ... n-1 instances and corresponding pointers if not a leaf child
	n2 := &Node{
		keys:     node.keys[half_n:],
		values:   node.values[half_n:],
		children: nil,
		parent:   node.parent,
		tree:     node.tree,
		isLeaf:   node.isLeaf,
	}

	if !node.isLeaf {
		n1.children = node.children[:half_n]
		n2.children = node.children[half_n:]
	}

	// Case 1: Node is root. Create new root with empty values and hook n1, n2.
	if node.tree.root == node {
		parent = &Node{
			keys:     []uint32{node.keys[half_n-1]},
			values:   []Addable{node.tree.aggregate.neutralElement, node.tree.aggregate.neutralElement},
			children: []*Node{n1, n2},
			parent:   nil,
			tree:     node.tree,
			isLeaf:   false,
		}
		n1.parent = parent
		n2.parent = parent
		parent.tree.root = parent
	} else {
		// Case 2: Node has parent. Let's insert n1, n2 and shift the keys, values and children to the right.
		parent = node.parent
		parent.keys = append(parent.keys, parent.keys[len(parent.keys)-1])
		parent.values = append(parent.values, parent.values[len(parent.values)-1])
		parent.children = append(parent.children, parent.children[len(parent.children)-1])

		for j, key := range parent.keys {
			if key >= node.keys[half_n-1] || j == int(parent.size()-1) {
				// change the following keys, values and children
				for i, _ := range parent.keys[j+1:] {
					parent.keys[len(parent.keys)-i-1] = parent.keys[len(parent.keys)-i-2]
					parent.values[len(parent.values)-i-2] = parent.values[len(parent.values)-i-3]
					parent.children[len(parent.children)-i-2] = parent.children[len(parent.children)-i-3]
				}
				// swap j to n1
				parent.children[j] = n1
				parent.keys[j] = node.keys[half_n-1]
				// the value at pos j stays the same

				// set n2, value of j+1 is already the value of j
				parent.children[j+1] = n2
				break
			}
		}
	}
	if parent.size()+1 > parent.tree.branchingFactor {
		parent.split()
	}
	return n1
}

func (node *Node) size() uint32 {
	return uint32(len(node.keys))
}
