package segmenttree

import "math"

type Node struct {
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

func (node *Node) size() uint32 {
	return uint32(len(node.keys))
}
