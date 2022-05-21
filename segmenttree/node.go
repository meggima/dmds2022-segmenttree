package segmenttree

import "math"

type Node struct {
	/*
			Size of a node is
			id: 4 byte (could be removed)
			keys: 4+4+4+ byte (reference, length, cap)  + b-1 * 4 byte
			values: assume float32 --> 12 byte + b * 4 byte
			children: 12 byte + l * 4 byte
			parent: 4 byte
			tree: 4 byte
			isLeaf: 1 byte

		This gives a total of  4 + 3*12 +(2*b-1) *4 + l * 4 + 8 +1 = 49 + (2*b-1) *4 + l * 4

		Go determines the memory allocation for our structs, it will pad bytes to make sure the final memory
		footprint is a multiple of 8 bytes
		let's assume the default page size of linux with 4 kb = 4096 bytes

		For simplicity set b = l, then we get a maximal branching factor b and a maximal leaf capacity l
		of b=l=337 to fit into one disk page of 4 kb (removing the id and the isLeaf flag would lead to 338).
	*/
	nodeId   uint32
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

		node.values[intervalIndex+2] = node.values[intervalIndex] // Original value // TODO breaks if empty node and len(values)<intervalIndex
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
		nodeId:   0, // TODO  This is basically useless, do we need it?
		keys:     node.keys[:half_n-1],
		values:   node.values[:half_n],
		children: nil,
		parent:   node.parent,
		tree:     node.tree,
		isLeaf:   node.isLeaf,
	}

	// N2 contains n/2 ... n-1 instances and corresponding pointers if not a leaf child
	n2 := &Node{
		nodeId:   0, // TODO  This is basically useless, do we need it?
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
			nodeId:   5, // This is basically useless
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

func (node *Node) imerge() {
	/*
		Merge two adjacent leaf intervals with equal aggregate values within one node.

		Following the practical advice in the paper and due to the overhead of the lookup we
			did not implement the case that two aggregate values of two neighbouring nodes could be combined.

		This procedure holds only if we never have more than two equal values beside each other. In practice this might be violated.
	*/
	if !node.isLeaf {
		return
	}
	for j, value := range node.values {
		if int(node.size()+1) > j && value == node.values[j+1] {
			node.keys = append(node.keys[:j], node.keys[j+1:]...)
			node.values = append(node.values[:j], node.values[j+1:]...)
			break
		} else if int(node.size()+1) == j && value == node.values[j+1] {
			node.keys = node.keys[:j]
			node.values = node.values[:j]
			break
		}
	}
}

func (node *Node) nmerge() {
	if node.size()+1 >= node.tree.branchingFactor/2 {
		// only nmerge if node is less than half full
		return
	}
	n := node.tree.branchingFactor
	halfN := int(math.Ceil(float64(n) / float64(2)))

	if !(node.size() != n) {
		panic("The node must hold exactly half_n_ceiled -1 elements. Thus one below the required minimum.")
	}

	if &node.tree.root == &node { // Case 1: node is root
		//node has only one child
		if len(node.children) == 1 {
			node.tree.root = node.children[0]
			node.tree.root.parent = nil
			for _, value := range node.tree.root.values {
				value = value.Add(node.values[0]) // TODO test this
			}
		}
		//do nothing
		return

	} else { // Case 2: node is not root
		// find the lef and right sibling
		var right_sibling *Node
		var left_sibling *Node
		var k int
		parent := node.parent
		for i, _ := range parent.children {
			if parent.children[i] == node {
				if i > 0 {
					left_sibling = parent.children[i-1]
				}
				if i <= int(parent.size()) {
					right_sibling = parent.children[i+1]
				}
				k = i
				break
			}
		}
		// Case2.1: If N' the right sibling of node has at least more than half_n +1 intervals, steal the first one  of N'! TODO test this
		if right_sibling != nil && int(right_sibling.size()) > halfN {
			for i, value := range node.values {
				node.values[i] = parent.values[k].Add(value)
			}
			parent.values[k] = node.tree.aggregate.neutralElement
			node.keys = append(node.keys, parent.keys[k])
			node.values = append(node.values, parent.values[k+1].Add(right_sibling.values[0]))
			if !node.isLeaf {
				node.children = append(node.children, right_sibling.children[0])
			}
			parent.keys[k] = right_sibling.keys[0]
			// cleanup sibling
			right_sibling.keys = right_sibling.keys[1:]
			right_sibling.values = right_sibling.values[1:]
			right_sibling.children = right_sibling.children[1:]

			return
		}

		// Case2.2: If N' the left sibling of N has more than half_n intervals Steal the last one of N'! TODO test this
		if left_sibling != nil && int(left_sibling.size()) > halfN {
			// in the paper N' the left sibling has now index k and N has index k+1 in the parent. Let's ignore this to keep things a bit more readable!
			for i, value := range node.values {
				node.values[i] = parent.values[k].Add(value)
			}
			parent.values[k] = node.tree.aggregate.neutralElement

			node.keys = append([]uint32{parent.keys[k-1]}, node.keys...)
			node.values = append([]Addable{parent.values[k-1].Add(left_sibling.values[len(left_sibling.values)-1])}, node.values...)
			if !node.isLeaf {
				node.children = append([]*Node{left_sibling.children[len(left_sibling.children)-1]}, node.children...)
			}
			parent.keys[k-1] = left_sibling.keys[len(left_sibling.keys)-2]
			// cleanup sibling
			left_sibling.keys = left_sibling.keys[:len(left_sibling.keys)-2]
			left_sibling.values = left_sibling.values[:len(left_sibling.values)-2]
			left_sibling.children = left_sibling.children[:len(left_sibling.children)-2]
			return
		}
		// Case2.3: Otherwise merge N with a sibling into a new node and place it in the parent of node.
		var n1 *Node
		var n2 *Node
		// TODO in practice there might be the case that left right sibling is nil. in this case we should take the left.
		if left_sibling != nil && left_sibling.size()+1 == node.tree.branchingFactor {
			n1 = left_sibling
			n2 = node
			k-- // so we know that k corresponds to n1
		} else if right_sibling != nil { // We need to loosen up the condition to make the example work. Removed  right_sibling.size()+1 == node.tree.branchingFactor
			n1 = node
			n2 = right_sibling
		} else {
			panic("no sibling has enough keys!")
		}

		newN := &Node{
			keys:     append(n1.keys, n2.keys...),
			values:   []Addable{},
			children: append(n1.children, n2.children...),
			parent:   parent,
			tree:     node.tree,
			isLeaf:   node.isLeaf,
		}

		newN.keys = append([]uint32{parent.keys[k]}, newN.keys...)

		for _, v := range n1.values {
			newN.values = append(newN.values, v.Add(parent.values[k]))
		}
		for _, v := range n2.values {
			newN.values = append(newN.values, v.Add(parent.values[k+1]))
		}
		// delete n1, n2 - this is not needed as we have a garbage collector
		parent.children[k] = newN
		parent.values[k] = node.tree.aggregate.neutralElement

		if int(parent.size()) > k {
			parent.keys = append(parent.keys[:k], parent.keys[k+1:]...)
			parent.values = append(parent.values[:k+1], parent.values[k+2:]...)
			if !node.isLeaf {
				parent.keys = append(parent.keys[:k+1], parent.keys[k+2:]...)
			}
		} else if int(parent.size()) == k {
			parent.keys = parent.keys[:k]
			parent.values = parent.values[:k+1]
			if !node.isLeaf {
				parent.keys = parent.keys[:k+1]
			}
		}
		// recurse: if the parent has now less then half_n nodes nmerge(parent)! TODO test this
		if int(parent.size())+1 <= halfN {
			parent.nmerge()
		}
	}
}
