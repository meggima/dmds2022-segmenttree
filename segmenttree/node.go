package segmenttree

import "math"

type Node struct {
	nodeId   uint32
	n        uint32
	keys     []uint32
	values   []Addable
	children []*Node
	parent   *Node
	tree     *SegmentTreeImpl
	isLeaf   bool
}

func (node *Node) findIntervalIndex(instant uint32) uint32 {
	var intervalIndex uint32 = 0

	for intervalIndex < node.n {
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
			return 0
		} else {
			childIndex := node.parent.findChildIndex(node)
			return node.parent.getIntervalStart(childIndex)
		}
	} else {
		return node.keys[index-1]
	}
}

func (node *Node) getIntervalEnd(index uint32) uint32 {
	if index == node.n {
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
	var intervals []Interval = make([]Interval, node.n+1)

	var i uint32 = 0
	for ; i <= node.n; i++ {
		start := node.getIntervalStart(i)
		end := node.getIntervalEnd(i)

		intervals[i] = NewInterval(start, end)
	}

	return intervals
}

func (node *Node) insert(intervalIndex int, tupleToInsert ValueIntervalTuple) {

}
